#!/usr/bin/env python3
import os
from pathlib import Path

BASE = Path("go-knowledge-base/01-Formal-Theory")

def write(path, content):
    full = BASE / path
    full.parent.mkdir(parents=True, exist_ok=True)
    with open(full, 'w', encoding='utf-8') as f:
        f.write(content)
    return full.stat().st_size

# File content templates
denotational = """# FT-011: Denotational Semantics
> **维度**: Formal Theory | **级别**: S

## 1. 形式化基础

### 1.1 指称语义定义

**定义 1.1 (指称语义)**
指称语义通过将程序映射到数学对象来描述程序含义:

$$[[ \\cdot ]] : \\text{Program} \\rightarrow \\text{Denotation}$$

**核心原则**:
- 合成性: 复合程序的指称由子程序指称组合而成
- 抽象性: 忽略实现细节，关注数学本质
- 精确性: 便于形式化证明

### 1.2 环境 (Environment)

**定义 1.2 (环境)**
环境是变量到值的映射:

$$\\rho \\in \\text{Env} = \\text{Var} \\rightharpoonup \\text{Val}$$

---

## 2. 简单表达式的指称语义

### 2.1 算术表达式

**定义 2.1 (常量指称)**
$$[[ n ]](\\rho) = n$$

**定义 2.2 (变量指称)**
$$[[ x ]](\\rho) = \\rho(x)$$

**定义 2.3 (加法指称)**
$$[[ e_1 + e_2 ]](\\rho) = [[ e_1 ]](\\rho) + [[ e_2 ]](\\rho)$$

### 2.2 命令式程序的指称

**定义 2.4 (存储指称)**
对于命令式程序，指称是存储转换器:
$$[[ c ]] : \\text{Env} \\rightarrow (\\text{Store} \\rightharpoonup \\text{Store}_{\\bot})$$

---

## 3. 不动点语义

### 3.1 循环的不动点定义

**定理 3.1 (While 循环的指称)**
While 循环的指称是以下函数的不动点:
$$F(f)(\\sigma) = \\begin{cases} f([[ c ]](\\rho)(\\sigma)) & \\text{if } [[ e ]](\\rho) = \\text{true} \\\\ \\sigma & \\text{if } [[ e ]](\\rho) = \\text{false} \\\\ \\bot & \\text{if } [[ e ]](\\rho) = \\bot \\end{cases}$$

$$[[ \\text{while } e \\text{ do } c ]](\\rho) = \\text{fix}(F)$$

### 3.2 Kleene 不动点定理

**定理 3.2 (Kleene 不动点定理)**
设 $f: D \\rightarrow D$ 是连续函数，$D$ 是有底元的 CPO，则:
$$\\text{fix}(f) = \\bigsqcup_{n \\geq 0} f^n(\\bot)$$
是最小不动点。

*证明*:
1. $f(\\text{fix}(f)) = \\text{fix}(f)$ (代入展开)
2. 若 $f(d) = d$，则 $\\bot \\sqsubseteq d$，归纳得 $f^n(\\bot) \\sqsubseteq d$
3. 因此 $\\text{fix}(f) \\sqsubseteq d$ (最小性) $\\square$

---

## 4. Go 语言的指称语义

### 4.1 函数类型的指称
**定义 4.1**
$$[[ \\text{func}(x_1 T_1) T_{ret} \\{ body \\} ]] = \\lambda v_1. [[ body ]](\\rho[x_1 \\mapsto v_1])$$

### 4.2 结构体的指称
**定义 4.2**
$$[[ t_S\\{e_1\\} ]](\\rho) = (t_S, \\{f_1: [[ e_1 ]](\\rho)\\})$$

### 4.3 Channel 的指称
**定义 4.3**
$$[[ \\text{chan } T ]] = \\text{Stream}([[ T ]])$$

---

## 5. 语义等价与证明

### 5.1 上下文等价

**定义 5.1 (完全抽象)**
语义是完全抽象的，如果:
$$e_1 \\cong e_2 \\iff [[ e_1 ]] = [[ e_2 ]]$$

**定理 5.1 (PCF 不完全抽象)**
Plotkin (1977): PCF 语言的标准指称语义**不是**完全抽象的。

---

## 6. 多元表征

### 6.1 语义风格对比

| 特性 | 操作语义 | 指称语义 | 公理语义 |
|------|----------|----------|----------|
| **粒度** | 执行步骤 | 数学函数 | 逻辑断言 |
| **中间状态** | 可见 | 隐藏 | 无 |
| **非终止** | 可表达 | 可表达 | 可表达 |
| **证明难度** | 中等 | 复杂 | 中等 |
| **典型应用** | 解释器 | 等价证明 | 验证 |

### 6.2 决策树
```
选择语义风格?
├── 需要数学抽象? → 指称语义
│   ├── 需要处理递归? → 域论 + 不动点
│   └── 需要完全抽象? → -game语义
└── 需要执行细节? → 操作语义
```

---

## 7. 实现示例

```go
package denotational

// Value 值类型
type Value interface {}

// Env 环境
type Env map[string]Value

// Denotation 指称函数
type Denotation func(Env) Value

// ConstDenotation 常量指称
func ConstDenotation(v Value) Denotation {
    return func(env Env) Value { return v }
}

// VarDenotation 变量指称  
func VarDenotation(name string) Denotation {
    return func(env Env) Value { return env[name] }
}

// PlusDenotation 加法指称
func PlusDenotation(left, right Denotation) Denotation {
    return func(env Env) Value {
        l := left(env).(int)
        r := right(env).(int)
        return l + r
    }
}
```

---

## 8. 参考文献

1. Winskel, G. (1993). *The Formal Semantics of Programming Languages*
2. Abramsky & Jung (1994). *Domain Theory*
3. Stoy (1977). *Denotational Semantics*

---

## 9. Toolkit

### 符号速查
| 符号 | 含义 |
|------|------|
| [[e]] | e 的指称 |
| rho | 环境 |
| bot | 底元/发散 |
| fix(f) | f 的最小不动点 |
| ⊔ | 最小上界 |

### 检查清单
- [ ] 定义值域 D 和偏序 ⊑
- [ ] 证明 D 是 CPO
- [ ] 证明语义函数连续
- [ ] 应用 Kleene 定理求不动点
- [ ] 验证完全抽象性
"""

axiomatic = """# FT-012: Axiomatic Semantics & Hoare Logic
> **维度**: Formal Theory | **级别**: S

## 1. 形式化基础

### 1.1 Hoare 三元组

**定义 1.1 (Hoare 三元组)**
$\\{P\\} C \\{Q\\}$ 表示: 如果前置条件 $P$ 成立，执行 $C$ 后 $Q$ 成立。

### 1.2 推理规则

**公理 1.1 (Skip)**
$$\\frac{\\}{\\{P\\} \\text{skip} \\{P\\}} \\text{(Skip)}$$

**公理 1.2 (赋值)**
$$\\frac{\\}{\\{Q[x := e]\\} x := e \\{Q\\}} \\text{(Assign)}$$

**规则 1.1 (顺序)**
$$\\frac{\\{P\\} C_1 \\{R\\} \\quad \\{R\\} C_2 \\{Q\\}}{\\{P\\} C_1; C_2 \\{Q\\}} \\text{(Seq)}$$

**规则 1.2 (条件)**
$$\\frac{\\{P \\land B\\} C_1 \\{Q\\} \\quad \\{P \\land \\neg B\\} C_2 \\{Q\\}}{\\{P\\} \\text{if } B \\text{ then } C_1 \\text{ else } C_2 \\{Q\\}} \\text{(If)}$$

**规则 1.3 (While)**
$$\\frac{\\{I \\land B\\} C \\{I\\}}{\\{I\\} \\text{while } B \\text{ do } C \\{I \\land \\neg B\\}} \\text{(While)}$$

---

## 2. 循环不变式

### 2.1 不变式的定义

**定义 2.1 (循环不变式)**
谓词 $I$ 是循环的不变式，如果:
$$\\{I \\land B\\} C \\{I\\}$$

### 2.2 寻找不变式

**定理 2.1 (求和循环)**
$$\\{n \\geq 0\\} \\text{Sum} \\{s = n(n+1)/2\\}$$

*证明*:
取 $I = (s = (i-1)i/2 \\land 1 \\leq i \\leq n+1)$
1. 初始化: $i=1, s=0$，$0 = 0$ ✓
2. 保持: $s' = s + i = i(i+1)/2$ ✓
3. 终止: $i = n+1$，$s = n(n+1)/2$ ✓

---

## 3. 最弱前置条件

### 3.1 Dijkstra 的谓词转换器

**定义 3.1 (wp)**
$\\text{wp}(C, Q)$ = 使得 $\\{P\\} C \\{Q\\}$ 成立的最弱 $P$

**公理 3.1 (wp-Skip)**
$$\\text{wp}(\\text{skip}, Q) = Q$$

**公理 3.2 (wp-赋值)**
$$\\text{wp}(x := e, Q) = Q[x := e]$$

**公理 3.3 (wp-顺序)**
$$\\text{wp}(C_1; C_2, Q) = \\text{wp}(C_1, \\text{wp}(C_2, Q))$$

---

## 4. Go 程序验证

### 4.1 Channel 公理

**公理 4.1 (发送)**
$$\\{P\\} ch <- v \\{P \\land \\text{sent}(ch, v)\\}$$

**公理 4.2 (接收)**
$$\\{P\\} x := <-ch \\{P \\land x = \\text{recv}(ch)\\}$$

### 4.2 Mutex 公理

**公理 4.3 (Lock/Unlock)**
$$\\{P\\} mu.Lock(); C; mu.Unlock() \\{Q\\}$$
其中 $C$ 访问共享变量。

---

## 5. 实现: 验证器

```go
package hoare

// Predicate 谓词
type Predicate func(Store) bool

// HoareTriple 三元组
type HoareTriple struct {
    Pre  Predicate
    Cmd  Command
    Post Predicate
}

// WP 计算最弱前置条件
func WP(cmd Command, post Predicate) Predicate {
    switch c := cmd.(type) {
    case Skip:
        return post
    case Assign:
        return func(s Store) bool {
            s2 := s.Copy()
            s2[c.Var] = eval(c.Expr, s)
            return post(s2)
        }
    case Seq:
        return WP(c.C1, WP(c.C2, post))
    default:
        return nil
    }
}

// Verify 验证三元组
func Verify(triple HoareTriple) bool {
    wp := WP(triple.Cmd, triple.Post)
    return implies(triple.Pre, wp)
}
```

---

## 6. 参考文献

1. Hoare, C.A.R. (1969). An Axiomatic Basis for Computer Programming
2. Dijkstra (1976). A Discipline of Programming
3. Gries (1981). The Science of Programming
"""

# Write files
print("Generating S-level documents...")
write("01-Semantics/02-Denotational-Semantics.md", denotational)
write("01-Semantics/03-Axiomatic-Semantics.md", axiomatic)
print("Done!")
