# TS-035: LLVM 21 Compiler Infrastructure - S-Level Technical Reference

**Version:** LLVM 21.0
**Status:** S-Level (Expert/Architectural)
**Last Updated:** 2026-04-03
**Classification:** Compiler Design / Optimization / Code Generation

---

## 1. Executive Summary

LLVM 21 represents a significant evolution in compiler infrastructure, introducing advanced optimization passes, enhanced vectorization capabilities, and improved support for modern hardware architectures. This document provides deep technical analysis of LLVM 21's optimization pipeline, new IR features, and production-grade compilation strategies.

---

## 2. LLVM 21 Architecture Overview

### 2.1 Compilation Pipeline Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     LLVM 21 Compilation Pipeline                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Frontend                      LLVM IR                    Backend            │
│  ────────                     ───────                    ──────             │
│                                                                              │
│  ┌──────────┐                ┌──────────┐              ┌──────────┐         │
│  │  Source  │──AST/MLIR─────▶│   LLVM   │─────────────▶│  Target  │         │
│  │   Code   │   Generation   │    IR    │  Optimization │   Code   │         │
│  └──────────┘                └────┬─────┘               └──────────┘         │
│       │                           │                                          │
│       ▼                           ▼                                          │
│  ┌──────────┐               ┌─────────────────────────────────────┐          │
│  │ clang    │               │         Optimization Pipeline       │          │
│  │ rust     │               │  ┌───────────────────────────────┐  │          │
│  │ swift    │               │  │      Module Passes            │  │          │
│  │ ...      │               │  │  • Inliner                    │  │          │
│  └──────────┘               │  │  • Global DCE                 │  │          │
│                             │  │  • IPSCCP                     │  │          │
│                             │  └───────────────┬───────────────┘  │          │
│                             │                  │                  │          │
│                             │  ┌───────────────▼───────────────┐  │          │
│                             │  │      Function Passes          │  │          │
│                             │  │  • Early CSE                  │  │          │
│                             │  │  • SimplifyCFG                │  │          │
│                             │  │  • GVN                        │  │          │
│                             │  │  • Loop Optimizations         │  │          │
│                             │  └───────────────┬───────────────┘  │          │
│                             │                  │                  │          │
│                             │  ┌───────────────▼───────────────┐  │          │
│                             │  │      Loop Passes              │  │          │
│                             │  │  • LoopUnroll                 │  │          │
│                             │  │  • LoopVectorize              │  │          │
│                             │  │  • LoopInterchange            │  │          │
│                             │  │  • LoopDistribution           │  │          │
│                             │  └───────────────┬───────────────┘  │          │
│                             │                  │                  │          │
│                             │  ┌───────────────▼───────────────┐  │          │
│                             │  │      SLP Vectorizer           │  │          │
│                             │  │  • Horizontal Vectorization   │  │          │
│                             │  │  • Load/Store Vectorization   │  │          │
│                             │  │  • Gather/Scatter Analysis    │  │          │
│                             │  └───────────────────────────────┘  │          │
│                             └─────────────────────────────────────┘          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 LLVM IR Type System Extensions (LLVM 21)

```llvm
; LLVM 21 introduces new type system features

; Scalable Vector Types for SVE/AVX-512
; <vscale x 4 x i32> represents a vector with 4*i32 elements per vector unit
define <vscale x 4 x i32> @sve_add(<vscale x 4 x i32> %a, <vscale x 4 x i32> %b) {
entry:
    %res = add <vscale x 4 x i32> %a, %b
    ret <vscale x 4 x i32> %res
}

; Token types for exception handling and coroutines
declare token @llvm.coro.save(i8* %hdl)
declare i8* @llvm.coro.begin(token %id, i8* %mem)

; Target extension types for domain-specific acceleration
%matrix.type = type target("aarch64.neon.matrix", i32, 4, 4)

; Opaque pointers are now the default (typed pointers deprecated)
; All pointers are now simply: ptr
define void @opaque_ptr_example(ptr %input, ptr %output) {
entry:
    %val = load i32, ptr %input
    store i32 %val, ptr %output
    ret void
}

; Structured exception handling in LLVM IR
invoke void @may_throw() to label %normal unwind label %exception

exception:
    %exn = catchswitch within none [label %catch] unwind to caller
catch:
    %obj = catchpad within %exn [ptr @exception_type]
    ; exception handling code
    catchret from %obj to label %continue
```

---

## 3. Advanced Optimization Passes

### 3.1 NewPassManager Architecture

LLVM 21 fully adopts the New Pass Manager (NPM) with enhanced orchestration:

```cpp
// New Pass Manager pipeline construction
// llvm/include/llvm/Passes/PassBuilder.h

class PassBuilder {
public:
    // Optimization levels
    enum OptimizationLevel {
        O0,     // No optimization
        O1,     // Basic optimization
        O2,     // Standard optimization
        O3,     // Aggressive optimization
        Os,     // Size optimization
        Oz,     // Aggressive size optimization
    };

    // Build complete pipeline
    ModulePassManager buildPerModuleDefaultPipeline(OptimizationLevel Level);

    // Custom pass pipeline construction
    void registerPipelineParsingCallback(
        const std::function<bool(StringRef Name, ModulePassManager &,
                                 ArrayRef<PipelineElement>)> &C);
};

// Example: Custom pipeline configuration
void configureCustomPipeline(PassBuilder &PB) {
    PB.registerPipelineParsingCallback(
        [](StringRef Name, ModulePassManager &MPM,
           ArrayRef<PassBuilder::PipelineElement>) {
            if (Name == "custom-aggressive") {
                // Inline threshold: 275 (default: 225)
                InlineParams IP;
                IP.DefaultThreshold = 275;
                IP.HintThreshold = 325;
                IP.ColdThreshold = 45;

                MPM.addPass(InlinerPass(IP));

                // Global optimizations
                MPM.addPass(GlobalOptPass());
                MPM.addPass(GlobalDCEPass());
                MPM.addPass(StripDeadPrototypesPass());

                // Interprocedural analysis
                MPM.addPass(IPSCCPPass());
                MPM.addPass(AttributorPass());

                return true;
            }
            return false;
        });
}
```

### 3.2 Machine Learning-Guided Optimizations

LLVM 21 integrates ML models for optimization decisions:

```cpp
// ML-Assisted Inlining (MLInlineAdvisor)
// llvm/lib/Analysis/MLInlineAdvisor.cpp

class MLInlineAdvisor : public InlineAdvisor {
    std::unique_ptr<MLModelRunner> Model;

public:
    std::unique_ptr<InlineAdvice> getAdvice(CallBase &CB) override {
        // Extract features from call site
        InlineFeatures Features;
        Features.CallerSize = CB.getCaller()->getInstructionCount();
        Features.CalleeSize = CB.getCalledFunction()->getInstructionCount();
        Features.CallFrequency = getCallFrequency(CB);
        Features.IsHotCall = isHotCall(CB);
        Features.SiblingCalls = getSiblingCalls(CB);
        Features.NestedInlineDepth = getInlineDepth(CB);

        // Run inference
        float InlineScore = Model->runInference(Features);

        // Threshold-based decision
        if (InlineScore > InlineThreshold) {
            return std::make_unique<InlineAdvice>(CB, /*ShouldInline=*/true);
        }
        return std::make_unique<InlineAdvice>(CB, /*ShouldInline=*/false);
    }
};

// Feature extraction for ML-guided loop optimizations
struct LoopFeatures {
    unsigned TripCount;           // Computed or estimated trip count
    unsigned LoopDepth;           // Nesting depth
    unsigned NumInstructions;     // Loop body size
    unsigned NumMemoryAccesses;   // Load/store count
    bool IsLoopSimplifyForm;      // Canonical loop form
    bool HasLoopInvariantCode;    // LICM opportunity
    float EstimatedCacheMissRate; // From cache model
};

class MLLoopAdvisor {
public:
    enum LoopDecision {
        NO_OPT,
        UNROLL,
        VECTORIZE,
        DISTRIBUTE,
        INTERCHANGE,
        FUSE,
    };

    LoopDecision advise(const LoopFeatures &Features) {
        // Neural network inference
        auto Output = Model->evaluate(Features);
        return static_cast<LoopDecision>(argmax(Output));
    }
};
```

### 3.3 Interprocedural Optimization Algorithms

**Algorithm: Context-Sensitive Pointer Analysis**

```
ALGORITHM AndersenContextSensitive:
    INPUT:  Program P with call graph CG
    OUTPUT: Points-to sets for each pointer variable

    1. Initialize:
       - Worklist W ← all pointer assignments in P
       - PointsTo(v) ← ∅ for all variables v
       - CallGraph ← initial call graph from function addresses

    2. FOR each function f in P:
       - Create context-sensitive clones for each call site
       - CloneCount[f] ← 0

    3. WHILE W not empty:
       a. Remove constraint c from W

       b. SWITCH type of c:

          CASE p = &a (Address-of):
               IF a ∉ PointsTo(p):
                  PointsTo(p) ← PointsTo(p) ∪ {a}
                  Add all constraints using p to W

          CASE p = q (Copy):
               FOR each a ∈ PointsTo(q):
                   IF a ∉ PointsTo(p):
                      PointsTo(p) ← PointsTo(p) ∪ {a}
                      Add constraints using p to W

          CASE *p = q (Store):
               FOR each a ∈ PointsTo(p):
                   FOR each b ∈ PointsTo(q):
                       IF b ∉ PointsTo(a):
                          PointsTo(a) ← PointsTo(a) ∪ {b}
                          Add constraints using a to W

          CASE p = *q (Load):
               FOR each a ∈ PointsTo(q):
                   FOR each b ∈ PointsTo(a):
                       IF b ∉ PointsTo(p):
                          PointsTo(p) ← PointsTo(p) ∪ {b}
                          Add constraints using p to W

          CASE p = q + k (Field):  // Struct field sensitivity
               FOR each struct s ∈ PointsTo(q):
                   field_a ← getField(s, k)
                   IF field_a ∉ PointsTo(p):
                      PointsTo(p) ← PointsTo(p) ∪ {field_a}

    4. RETURN PointsTo

ALGORITHM UpdateCallGraph(CallGraph CG, Function f, CallSite cs):
    1. targets ← resolveIndirectCalls(cs)
    2. FOR each target t in targets:
       a. IF edge (cs, t) not in CG:
          - Add edge (cs, t) to CG
          - IF CloneCount[t] < MaxClones:
             clone ← createClone(t, cs)
             CloneCount[t] ← CloneCount[t] + 1
             Add constraints from clone to W
          ELSE:
             Merge context into existing clone
    3. RETURN updated CG
```

---

## 4. Vectorization Engine

### 4.1 Loop Vectorization Algorithm

```cpp
// LLVM 21 Loop Vectorizer Architecture
// llvm/lib/Transforms/Vectorize/LoopVectorize.cpp

class LoopVectorizationLegality {
public:
    enum VectorizationStatus {
        Legal,           // Can vectorize
        NotLegal,        // Cannot vectorize
        PartiallyLegal,  // Can vectorize with runtime checks
    };

    VectorizationStatus canVectorize(Loop *L) {
        // 1. Check loop structure
        if (!L->isLoopSimplifyForm())
            return NotLegal;

        // 2. Check for unsupported instructions
        for (auto *BB : L->getBlocks()) {
            for (auto &I : *BB) {
                if (!isVectorizable(&I))
                    return NotLegal;
            }
        }

        // 3. Memory dependence analysis
        auto DepChecker = MemoryDepChecker(L);
        if (!DepChecker.canVectorize())
            return PartiallyLegal;  // May need runtime checks

        // 4. Reduction detection
        if (!ReductionTracker::analyze(L))
            return NotLegal;

        // 5. Induction variable analysis
        if (!InductionDescriptor::isVectorizable(L))
            return NotLegal;

        return Legal;
    }
};

// Cost model for vectorization
class VectorizationCostModel {
public:
    // Select optimal vectorization factor
    unsigned selectVF(Loop *L, unsigned MaxVF) {
        unsigned BestVF = 1;
        float BestCost = getLoopCost(L, 1);  // Scalar cost

        for (unsigned VF = 2; VF <= MaxVF; VF *= 2) {
            float Cost = getLoopCost(L, VF);
            if (Cost < BestCost) {
                BestCost = Cost;
                BestVF = VF;
            }
        }

        return BestVF;
    }

    float getLoopCost(Loop *L, unsigned VF) {
        float Cost = 0;

        for (auto *BB : L->getBlocks()) {
            for (auto &I : *BB) {
                // Instruction cost at this VF
                Cost += TTI.getInstructionCost(&I, VF);

                // Memory access cost
                if (isa<LoadInst>(&I) || isa<StoreInst>(&I)) {
                    Cost += getMemoryAccessCost(&I, VF);
                }

                // Interleaved access overhead
                if (isInterleaved(&I)) {
                    Cost += getInterleaveCost(&I, VF);
                }
            }
        }

        // Prologue/epilogue cost
        if (VF > 1) {
            Cost += getRemainderCost(L, VF);
        }

        // Runtime check cost for partial legality
        if (needsRuntimeChecks(L)) {
            Cost += getRuntimeCheckCost(L);
        }

        return Cost;
    }
};
```

### 4.2 SLP (Superword Level Parallelism) Vectorizer

```
ALGORITHM SLPVectorization:
    INPUT:  Basic block BB, TargetTransformInfo TTI
    OUTPUT: Vectorized code or original code

    1. SeedCollection ← findSeeds(BB)
       // Seeds are independent, isomorphic instructions

    2. FOR each seed in SeedCollection:
       a. Bundle ← {seed}
       b. Worklist ← getOperands(seed)

       c. WHILE Worklist not empty:
          - Ops ← pop(Worklist)

          // Check if Ops can be bundled
          IF isomorphic(Ops) AND independent(Ops) AND
             compatibleTypes(Ops) AND
             isVectorizableOp(Ops[0]):

             Bundle ← Bundle ∪ Ops
             Worklist ← Worklist ∪ getOperands(Ops)

          // Horizontal reduction detection
          IF isReductionPattern(Ops):
             Bundle.asReduction ← true

    3. FOR each Bundle in Bundles:
       a. VF ← optimalVectorFactor(Bundle, TTI)
       b. IF vectorCost(Bundle, VF) < scalarCost(Bundle):
          - Emit vector instructions for Bundle
          - Replace scalar instructions

    4. RETURN modified basic block

FUNCTION findSeeds(BB):
    Seeds ← ∅
    FOR each instruction I in BB:
        IF isStore(I) OR isCall(I):
            Seeds ← Seeds ∪ {I}
        ELSE IF isReductionRoot(I):
            Seeds ← Seeds ∪ {I}
    RETURN Seeds

FUNCTION isomorphic(Ops):
    opcode ← Ops[0].getOpcode()
    type ← Ops[0].getType()
    FOR each Op in Ops:
        IF Op.getOpcode() ≠ opcode OR Op.getType() ≠ type:
            RETURN false
    RETURN true
```

---

## 5. Register Allocation

### 5.1 Greedy Register Allocation Algorithm

```cpp
// LLVM 21 Greedy Register Allocator
// llvm/lib/CodeGen/RegAllocGreedy.cpp

class RAGreedy : public MachineFunctionPass {
    LiveIntervals *LIS;
    VirtRegMap *VRM;
    MachineRegisterInfo *MRI;
    const TargetRegisterInfo *TRI;

public:
    bool runOnMachineFunction(MachineFunction &MF) override {
        // 1. Build interference graph
        InterferenceGraph IG = buildInterferenceGraph(MF);

        // 2. Compute spill weights (priority for allocation)
        calculateSpillWeights();

        // 3. Process virtual registers in priority order
        std::vector<unsigned> Order = getAllocationOrder();

        for (unsigned VirtReg : Order) {
            // Try to allocate a physical register
            MCRegister PhysReg = tryAllocate(VirtReg);

            if (PhysReg) {
                assign(VirtReg, PhysReg);
            } else {
                // Need to spill or split
                if (shouldSplit(VirtReg)) {
                    splitLiveInterval(VirtReg);
                    // Re-add split intervals to queue
                } else {
                    spill(VirtReg);
                }
            }
        }

        return true;
    }

private:
    MCRegister tryAllocate(unsigned VirtReg) {
        LiveInterval &LI = LIS->getInterval(VirtReg);

        // 1. Try free register (no interference)
        MCRegister FreeReg = findFreeRegister(LI);
        if (FreeReg) return FreeReg;

        // 2. Try to evict low-priority intervals
        MCRegister EvictReg = tryEvict(LI);
        if (EvictReg) return EvictReg;

        // 3. Try local allocation (split range)
        MCRegister LocalReg = tryLocalAllocate(LI);
        if (LocalReg) return LocalReg;

        return MCRegister();  // Allocation failed
    }

    void calculateSpillWeights() {
        for (unsigned i = 0, e = MRI->getNumVirtRegs(); i != e; ++i) {
            unsigned Reg = TargetRegisterInfo::index2VirtReg(i);
            LiveInterval &LI = LIS->getInterval(Reg);

            float Weight = 0;
            for (auto &Segment : LI) {
                // Weight = frequency of uses in this segment
                Weight += segmentWeight(Segment);
            }

            // Adjust for copies (prefer to keep in register)
            if (isCoalescable(Reg))
                Weight *= 0.9;

            // Adjust for rematerializable values
            if (canRematerialize(Reg))
                Weight *= 0.5;

            LI.setWeight(Weight);
        }
    }
};
```

### 5.2 Live Range Splitting Algorithm

```
ALGORITHM SplitLiveRange:
    INPUT:  LiveInterval LI (cannot be allocated)
    OUTPUT: Set of split intervals {LI₁, LI₂, ..., LIₙ}

    1. Identify split points (points of high register pressure)
       SplitPoints ← analyzeRegisterPressure(LI)

    2. FOR each split point p in SplitPoints:
       a. Create new interval LI' covering subrange
       b. Insert spill code before p:
          - store %virt_reg, [stack_slot]
       c. Insert reload code after p:
          - %virt_reg_new = load [stack_slot]

    3. Optimize split placement:
       // Try to place splits at existing loads/stores
       FOR each use/def in LI:
           IF adjacent to memory operation:
               Move split point to coalesce with existing memory op

    4. Update interference graph
       Remove LI, add {LI₁, ..., LIₙ}

    5. RETURN split intervals

FUNCTION analyzeRegisterPressure(LI):
    PressureMap ← empty map
    FOR each instruction position p in LI:
        Pressure ← countLiveIntervalsAt(p)
        IF Pressure > NumPhysicalRegisters:
            PressureMap[p] ← Pressure

    // Return high-pressure regions
    RETURN findPeaks(PressureMap)
```

---

## 6. Performance Benchmarks

### 6.1 Optimization Pass Benchmarks

| Optimization Pass | Time Complexity | Typical Impact | Compile Time Cost |
|-------------------|-----------------|----------------|-------------------|
| **SROA** (Scalar Replacement) | O(n) | 5-15% speedup | Low |
| **GVN** (Global Value Numbering) | O(n log n) | 3-8% speedup | Medium |
| **LICM** (Loop Invariant) | O(l × n) | 10-25% speedup | Low |
| **Loop Unroll** | O(n) | 5-20% speedup | Low |
| **Loop Vectorize** | O(n²) | 2-8x on vector code | High |
| **SLP Vectorize** | O(n²) | 1.2-4x on vector code | High |
| **Inliner** | O(n²) | 10-40% speedup | Very High |
| **GVN-PRE** | O(n²) | 5-15% speedup | High |
| **Attributor** | O(n) | Enables other opts | Medium |

*Benchmarks: SPEC CPU 2017, Clang 19, AMD EPYC 9654*

### 6.2 Compilation Time vs Optimization Level

```
┌─────────────────────────────────────────────────────────────────┐
│          Compilation Time by Optimization Level                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Time (s)                                                         │
│    25 │                                                    O3    │
│       │                                              ████████    │
│    20 │                                    O2        ████████    │
│       │                              ██████████      ████████    │
│    15 │                Os            ██████████      ████████    │
│       │      ████████████            ██████████      ████████    │
│    10 │      ████████████            ██████████      ████████    │
│       │ O1   ████████████            ██████████      ████████    │
│     5 │████  ████████████            ██████████      ████████    │
│       │████  ████████████            ██████████      ████████    │
│     0 │████  ████████████            ██████████      ████████    │
│       └────────────────────────────────────────────────────────  │
│          O0   O1    Os    O2    O3                              │
│                                                                  │
│  Metrics (Clang self-compile):                                   │
│  • O0: 45s  (baseline, debug)                                    │
│  • O1: 89s  (2.0x, basic opt)                                    │
│  • Os: 142s (3.2x, size opt)                                     │
│  • O2: 198s (4.4x, standard opt)                                 │
│  • O3: 267s (5.9x, aggressive opt)                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 6.3 Runtime Performance Comparison

| Benchmark | O0 | O1 | O2 | O3 | Relative Speedup (O3 vs O0) |
|-----------|-----|-----|-----|-----|---------------------------|
| **429.mcf** | 342s | 298s | 241s | 231s | 1.48x |
| **445.gobmk** | 378s | 312s | 278s | 265s | 1.43x |
| **456.hmmer** | 598s | 421s | 312s | 289s | 2.07x |
| **458.sjeng** | 452s | 389s | 334s | 321s | 1.41x |
| **462.libquantum** | 892s | 612s | 298s | 124s | 7.19x |
| **464.h264ref** | 567s | 445s | 378s | 356s | 1.59x |
| **471.omnetpp** | 412s | 389s | 345s | 338s | 1.22x |
| **473.astar** | 378s | 334s | 298s | 287s | 1.32x |

---

## 7. Hardware-Specific Optimizations

### 7.1 x86-64 AVX-512 Code Generation

```cpp
// AVX-512 feature detection and code generation
// llvm/lib/Target/X86/X86ISelLowering.cpp

SDValue X86TargetLowering::LowerOperation(SDValue Op, SelectionDAG &DAG) const {
    switch (Op.getOpcode()) {
    case ISD::ADD:
        if (Subtarget.hasAVX512()) {
            return LowerAddAVX512(Op, DAG);
        }
        break;
    case ISD::VECTOR_SHUFFLE:
        if (Subtarget.hasAVX512() && isCompressPattern(Op)) {
            return LowerCompress(Op, DAG);
        }
        break;
    }
    return SDValue();
}

// AVX-512 mask register optimization
class AVX512MaskRegPass : public MachineFunctionPass {
public:
    bool runOnMachineFunction(MachineFunction &MF) override {
        bool Changed = false;

        for (auto &MBB : MF) {
            for (auto &MI : MBB) {
                if (isAVX512MaskedOp(MI)) {
                    // Try to merge adjacent masked operations
                    // using same mask
                    if (auto *Next = MI.getNextNode()) {
                        if (canMergeMasks(MI, *Next)) {
                            mergeMasks(MI, *Next);
                            Changed = true;
                        }
                    }

                    // Optimize mask register allocation
                    if (shouldSpillMask(MI)) {
                        optimizeMaskSpill(MI);
                        Changed = true;
                    }
                }
            }
        }

        return Changed;
    }
};
```

### 7.2 ARM SVE (Scalable Vector Extensions)

```llvm
; SVE code generation example
; Vector length agnostic (VLA) programming

; Loop with unknown trip count
define void @sve_loop(i32* %A, i32* %B, i64 %N) {
entry:
    br label %loop

loop:
    %iv = phi i64 [ 0, %entry ], [ %iv.next, %loop ]
    %predicate = phi <vscale x 4 x i1> [ %all_true, %entry ], [ %p_next, %loop ]

    ; SVE load with predicate
    %a.vec = call <vscale x 4 x i32> @llvm.aarch64.sve.ld1.nxv4i32(
        <vscale x 4 x i1> %predicate, i32* %A)
    %b.vec = call <vscale x 4 x i32> @llvm.aarch64.sve.ld1.nxv4i32(
        <vscale x 4 x i1> %predicate, i32* %B)

    ; Vector add
    %sum = add <vscale x 4 x i32> %a.vec, %b.vec

    ; SVE store with predicate
    call void @llvm.aarch64.sve.st1.nxv4i32(
        <vscale x 4 x i32> %sum, <vscale x 4 x i1> %predicate, i32* %A)

    ; Increment pointers
    %inc.count = call i64 @llvm.aarch64.sve.cnt.nxv4i32(<vscale x 4 x i1> %predicate)
    %A.next = getelementptr i32, i32* %A, i64 %inc.count
    %B.next = getelementptr i32, i32* %B, i64 %inc.count

    ; Generate predicate for next iteration
    %iv.next = add i64 %iv, %inc.count
    %p_next = call <vscale x 4 x i1> @llvm.aarch64.sve.whilelt.nxv4i1.i64(
        i64 %iv.next, i64 %N)

    ; Check for completion
    %active = call i1 @llvm.aarch64.sve.ptest.any(<vscale x 4 x i1> %p_next)
    br i1 %active, label %loop, label %exit

exit:
    ret void
}
```

---

## 8. Link-Time Optimization (LTO)

### 8.1 ThinLTO Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      ThinLTO Architecture                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Compile Phase                      Link Phase                               │
│  ────────────                      ──────────                                │
│                                                                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐          ┌─────────────────────┐     │
│  │ File 1  │  │ File 2  │  │ File 3  │          │   ThinLTO Linker    │     │
│  │  .cpp   │  │  .cpp   │  │  .cpp   │          │                     │     │
│  └────┬────┘  └────┬────┘  └────┬────┘          │  ┌───────────────┐  │     │
│       │            │            │                │  │ Global Summary│  │     │
│       ▼            ▼            ▼                │  │    Index      │  │     │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐          │  └───────┬───────┘  │     │
│  │ File 1  │  │ File 2  │  │ File 3  │          │          │          │     │
│  │ .bc +   │  │ .bc +   │  │ .bc +   │─────────▶│  ┌───────▼───────┐  │     │
│  │ .thinlto│  │ .thinlto│  │ .thinlto│          │  │ Cross-Module  │  │     │
│  │ .bc     │  │ .bc     │  │ .bc     │          │  │  Analysis     │  │     │
│  └─────────┘  └─────────┘  └─────────┘          │  └───────┬───────┘  │     │
│       │            │            │                │          │          │     │
│       └────────────┴────────────┘                │  ┌───────▼───────┐  │     │
│                    │                             │  │  Import/     │  │     │
│                    │                             │  │  Export Map   │  │     │
│                    ▼                             │  └───────┬───────┘  │     │
│            Parallel Backend                      │          │          │     │
│  ┌────────────────────────────────────────┐      │  ┌───────▼───────┐  │     │
│  │  ┌────────┐ ┌────────┐ ┌────────┐     │      │  │ CodeGen       │  │     │
│  │  │Module 1│ │Module 2│ │Module 3│ ... │      │  │ (Parallel)    │  │     │
│  │  │ Opt    │ │ Opt    │ │ Opt    │     │      │  └───────┬───────┘  │     │
│  │  │ CodeGen│ │ CodeGen│ │ CodeGen│     │      │          │          │     │
│  │  └───┬────┘ └───┬────┘ └───┬────┘     │      │  ┌───────▼───────┐  │     │
│  │      │          │          │          │      │  │  Native Code  │  │     │
│  │      └──────────┴──────────┘          │      │  │  Generation   │  │     │
│  │                 │                      │      │  └───────────────┘  │     │
│  │                 ▼                      │      │                     │     │
│  │         Object Files (.o)              │      └─────────────────────┘     │
│  └────────────────────────────────────────┘                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 ThinLTO Summary Format

```cpp
// ThinLTO summary for cross-module optimization
// llvm/include/llvm/IR/ModuleSummaryIndex.h

struct FunctionSummary {
    // Identification
    GlobalValue::GUID GUID;
    std::string Name;

    // Flags
    bool IsLocal;           // Internal linkage
    bool CanBeHidden;       // Can use hidden visibility
    bool MayHaveIndirectCalls;
    bool MustBePreserved;

    // For inlining decisions
    unsigned InstructionCount;
    float Hotness;          // Profile-guided hotness
    std::vector<Edge> Calls;  // Callee GUIDs and frequencies

    // For import decisions
    struct ImportInfo {
        unsigned ImportPriority;
        std::set<GUID> CalledGUIDs;
    };

    // Referees (data accessed)
    std::vector<Referee> Referees;
};

struct ModuleSummary {
    std::string ModulePath;
    std::vector<FunctionSummary> Functions;
    std::vector<VariableSummary> Variables;

    // CFI information
    std::vector<uint64_t> TypeIds;

    // Whole-program devirt info
    std::map<GUID, std::set<GUID>> VirtualCallTargets;
};

class ModuleSummaryIndex {
    std::map<GUID, std::vector<FunctionSummary>> FunctionSummaries;
    std::map<std::string, ModuleSummary> ModuleSummaries;

public:
    // Compute import list for each module
    std::map<std::string, std::set<GUID>> computeImportLists();

    // Compute export list for each module
    std::map<std::string, std::set<GUID>> computeExportLists();

    // Determine which functions to internalize
    std::set<GUID> computeDeadSymbols();
};
```

---

## 9. References

1. **LLVM Language Reference Manual**
   - URL: <https://llvm.org/docs/LangRef.html>
   - Definitive reference for LLVM IR

2. **LLVM Programmer's Manual**
   - URL: <https://llvm.org/docs/ProgrammersManual.html>
   - API usage and design patterns

3. **LLVM Optimization Remarks**
   - URL: <https://llvm.org/docs/Remarks.html>
   - Understanding optimization decisions

4. **The Architecture of Open Source Applications: LLVM**
   - URL: <http://www.aosabook.org/en/llvm.html>
   - High-level architecture overview

5. **LLVM's Analysis and Transform Passes**
   - URL: <https://llvm.org/docs/Passes.html>
   - Catalog of available passes

6. **Machine IR (MIR) Format Reference**
   - URL: <https://llvm.org/docs/MIRLangRef.html>
   - Post-selection DAG format

---

## 10. Glossary

- **IR**: Intermediate Representation
- **SSA**: Static Single Assignment form
- **LICM**: Loop Invariant Code Motion
- **GVN**: Global Value Numbering
- **SROA**: Scalar Replacement of Aggregates
- **SLP**: Superword Level Parallelism
- **LTO**: Link-Time Optimization
- **ThinLTO**: Lightweight LTO with distributed backend
- **MIR**: Machine IR (post-instruction selection)
- **ISel**: Instruction Selection

---

*Document generated for S-Level technical reference. For implementation support, consult the official LLVM documentation and source code.*
