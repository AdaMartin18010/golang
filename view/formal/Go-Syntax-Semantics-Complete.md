# Go语言完整语法语义形式化 (Featherweight Go / FG & FGG)

## 思维导图：Go语言形式化体系

```
                    Go语言形式化体系
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
    Featherweight      Go 1.22规范        Go泛型(FGG)
       Go (FG)                              │
        │                  │               │
    ┌───┴───┐         ┌────┴────┐     ┌────┴────┐
    ▼       ▼         ▼         ▼     ▼         ▼
  结构    接口      词法       语法   类型参数   单态化
  子类型  方法集    语法      语义    约束      翻译
```

---

## 1. 完整EBNF语法 (Go 1.22)

### 1.1 词法元素

**BNF (Lexical)**:

```
<SourceFile>      ::= <PackageClause> ";" { <ImportDecl> ";" } { <TopLevelDecl> ";" }

<PackageClause>   ::= "package" <PackageName>
<PackageName>     ::= <identifier>

<ImportDecl>      ::= "import" ( <ImportSpec> | "(" { <ImportSpec> ";" } ")" )
<ImportSpec>      ::= [ "." | <PackageName> ] <ImportPath>
<ImportPath>      ::= <string_lit>

<TopLevelDecl>    ::= <Declaration> | <FunctionDecl> | <MethodDecl>
<Declaration>     ::= <ConstDecl> | <TypeDecl> | <VarDecl>

<identifier>      ::= <letter> { <letter> | <unicode_digit> }
<letter>          ::= <unicode_letter> | "_"
```

### 1.2 类型系统语法

```
<Type>            ::= <TypeName> | <TypeLit> | "(" <Type> ")"
<TypeName>        ::= <identifier> | <QualifiedIdent>
<QualifiedIdent>  ::= <PackageName> "." <identifier>

<TypeLit>         ::= <ArrayType> | <StructType> | <PointerType>
                    | <FunctionType> | <InterfaceType> | <SliceType>
                    | <MapType> | <ChannelType>

<ArrayType>       ::= "[" <ArrayLength> "]" <ElementType>
<ArrayLength>     ::= <Expression>
<ElementType>     ::= <Type>

<SliceType>       ::= "[" "]" <ElementType>

<StructType>      ::= "struct" "{" { <FieldDecl> ";" } "}"
<FieldDecl>       ::= (<FieldNameList> <Type> | <EmbeddedField>) [ <Tag> ]
<FieldNameList>   ::= <FieldName> { "," <FieldName> }
<FieldName>       ::= <identifier>
<EmbeddedField>   ::= [ "*" ] <TypeName>
<Tag>             ::= <string_lit>

<PointerType>     ::= "*" <BaseType>
<BaseType>        ::= <Type>

<FunctionType>    ::= "func" <Signature>
<Signature>       ::= <Parameters> [ <Result> ]
<Result>          ::= <Parameters> | <Type>
<Parameters>      ::= "(" [ <ParameterList> [ "," ] ] ")"
<ParameterList>   ::= <ParameterDecl> { "," <ParameterDecl> }
<ParameterDecl>   ::= [ <IdentifierList> ] [ "..." ] <Type>
<IdentifierList>  ::= <identifier> { "," <identifier> }

<InterfaceType>   ::= "interface" "{" { <InterfaceElem> ";" } "}"
<InterfaceElem>   ::= <MethodElem> | <TypeElem>
<MethodElem>      ::= <MethodName> <Signature>
<MethodName>      ::= <identifier>
<TypeElem>        ::= <TypeTerm> { "|" <TypeTerm> }
<TypeTerm>        ::= <Type> | UnderlyingType
<UnderlyingType>  ::= "~" <Type>
```

### 1.3 语句与表达式

```
<Statement>       ::= <Declaration> | <LabeledStmt> | <SimpleStmt>
                    | <GoStmt> | <ReturnStmt> | <BreakStmt>
                    | <ContinueStmt> | <GotoStmt> | <FallthroughStmt>
                    | <Block> | <IfStmt> | <SwitchStmt> | <SelectStmt>
                    | <ForStmt> | <DeferStmt>

<SimpleStmt>      ::= <EmptyStmt> | <ExpressionStmt> | <SendStmt>
                    | <IncDecStmt> | <Assignment> | <ShortVarDecl>

<GoStmt>          ::= "go" <Expression>
<ChannelType>     ::= ( "chan" | "chan" "<-" | "<-" "chan" ) <ElementType>
<SendStmt>        ::= <Channel> "<-" <Expression>
<Channel>         ::= <Expression>

<RecvStmt>        ::= [ <ExpressionList> "=" ] <RecvExpr>
<RecvExpr>        ::= <- <Expression>

<ForStmt>         ::= "for" [ <Condition> | <ForClause> | <RangeClause> ] <Block>
<Condition>       ::= <Expression>
<ForClause>       ::= [ <InitStmt> ] ";" [ <Condition> ] ";" [ <PostStmt> ]
<InitStmt>        ::= <SimpleStmt>
<PostStmt>        ::= <SimpleStmt>
<RangeClause>     ::= [ <ExpressionList> "=" | <IdentifierList> ":=" ] "range" <Expression>
```

---

## 2. Featherweight Go (FG) 核心演算

### 2.1 FG 抽象语法

**定义 2.1.1** (FG语法).
$$
\begin{align}
\text{程序 } P &::= \overline{D} \; \text{func main}() \{ \_ = e \} \\
\text{声明 } D &::= \text{type } t_S \; \text{struct}\{\overline{f \; \tau}\} \\
&\quad \mid \text{type } t_I \; \text{interface}\{\overline{S}\} \\
&\quad \mid \text{func } (x \; t_S) \; m(\overline{x \; \tau}) \; \tau \; \{ \text{return } e \} \\
\text{表达式 } e &::= x \mid e.m(\overline{e}) \mid t_S\{\overline{e}\} \mid e.f \mid e.(\tau) \\
\text{类型 } \tau &::= t_S \mid t_I \\
\text{方法规约 } S &::= m(\overline{\tau}) \; \tau
\end{align}
$$

### 2.2 FG 类型系统

**规则 2.2.1** (结构类型).
结构类型 **t_S** 包含字段 **f̄**。

**规则 2.2.2** (接口满足).
结构类型 **t_S** 满足接口 **t_I** 当且仅当 **t_S** 实现了 **t_I** 的所有方法：
$$
\frac{\forall m(\overline{\tau})\tau \in t_I. \; t_S \text{ has } m(\overline{\tau})\tau}{t_S <: t_I}
$$

**规则 2.2.3** (方法调用).
$$
\frac{\Gamma \vdash e : t_S \quad \text{method}(t_S, m) = (\overline{x \; \tau}) \; \tau_r \; \{ \text{return } e' \} \quad \Gamma \vdash \overline{e} : \overline{\tau}}{\Gamma \vdash e.m(\overline{e}) : \tau_r}
$$

### 2.3 FG 操作语义 (小步)

**规则 2.3.1** (字段选择).
$$
t_S\{\overline{v}\}.f_i \longrightarrow v_i
$$

**规则 2.3.2** (方法调用).
$$
t_S\{\overline{v}\}.m(\overline{v'}) \longrightarrow e[x \mapsto t_S\{\overline{v}\}, \overline{x_i \mapsto v'_i}]
$$
其中 **method(t_S, m) = func(x t_S) m(𝑥̄ τ) τ { return e }**。

**规则 2.3.3** (类型断言 - 成功).
$$
t_S\{\overline{v}\}.(t_S) \longrightarrow t_S\{\overline{v}\}
$$

**规则 2.3.4** (类型断言 - 失败).
$$
t_S\{\overline{v}\}.(t'_I) \longrightarrow \text{panic} \quad \text{if } t_S \not<: t'_I
$$

---

## 3. Featherweight Generic Go (FGG)

### 3.1 FGG 抽象语法

**定义 3.1.1** (FGG语法).
$$
\begin{align}
\text{类型形参 } \Phi &::= \overline{\alpha \; \gamma} \\
\text{类型 } \tau &::= \alpha \mid t(\overline{\tau}) \\
\text{约束 } \gamma &::= \text{interface}\{\overline{S}\} \\
\text{声明 } D &::= \text{type } t[\Phi] \; \text{struct}\{\overline{f \; \tau}\} \\
&\quad \mid \text{type } t[\Phi] \; \text{interface}\{\overline{S}\} \\
&\quad \mid \text{func } (x \; t_S[\overline{\tau}]) \; m[\Psi](\overline{x \; \tau}) \; \tau \; \{ \text{return } e \}
\end{align}
$$

### 3.2 FGG 单态化翻译

**定义 3.2.1** (单态化).
单态化将 FGG 程序翻译为 FG 程序：
$$
\llbracket P \rrbracket_{FGG \to FG} = P^\dagger
$$

**规则 3.2.2** (类型实例化).
$$
\llbracket t[\overline{\tau}] \rrbracket_\eta = t\langle\overline{\llbracket \tau \rrbracket_\eta}\rangle
$$

**规则 3.2.3** (方法实例化).
$$
\llbracket e.m[\overline{\tau}](\overline{e}) \rrbracket_\eta = \llbracket e \rrbracket_\eta . m^\dagger(\overline{\llbracket \tau \rrbracket_\eta})(\overline{\llbracket e \rrbracket_\eta})
$$

---

## 4. Go内存模型 (补充)

见 `Go-Memory-Model-Formalization.md`。

---

## 参考文献

1. Griesemer, R. et al. (2020). *Featherweight Go*. OOPSLA.
2. Go Authors. (2024). *The Go Programming Language Specification*.
3. Ellis, S. & Zhu, S. (2022). *Featherweight Go Formalization*.

---

*文档版本: 2026-03-29 | 形式化等级: 完整EBNF + FG/FGG演算*
