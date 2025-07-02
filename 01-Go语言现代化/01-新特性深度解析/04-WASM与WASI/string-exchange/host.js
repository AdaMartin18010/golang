/*
    ======================================================================
    Node.js Host for Go WASM String Exchange Example
    ======================================================================

    🎯 目的:
    演示一个 Node.js 宿主环境如何加载 Go WASM 模块，并通过管理
    WASM 内存来与 Go 函数交换字符串数据。

    ⚙️ 如何准备:
    1. 确保你已经编译了 Go 源代码:
       $ go build -o main.wasm .
    2. 确保你的 Node.js 版本支持 WASI (v16+ 推荐)。

    🚀 如何运行:
    $ node host.js

    🔍 预期输出:
    > Go WASM module with string support loaded.
    > Calling Go's greet function with input: "World"
    > Received response from Go: "Hello, World!"
*/

import { readFile } from 'node:fs/promises';
import { WASI } from 'wasi';

async function run() {
    // 1. 初始化 WASI 环境
    const wasi = new WASI({
        version: 'preview1',
        // 允许 WASM 模块访问标准输入输出
        args: process.argv,
        env: process.env,
        preopens: {
            '/sandbox': './',
        },
    });

    // 2. 读取并编译 WASM 模块
    const wasmPath = './main.wasm';
    const wasmBytes = await readFile(wasmPath);
    const module = await WebAssembly.compile(wasmBytes);

    // 3. 实例化 WASM 模块，并传入 WASI 的导入对象
    const instance = await WebAssembly.instantiate(module, wasi.getImportObject());
    wasi.start(instance);

    // 4. 从实例中获取导出的函数和内存
    const { greet, malloc, free } = instance.exports;
    const memory = instance.exports.memory;
    
    // --- 将字符串从 JS 传递到 Go ---

    const inputString = "World";
    console.log(`Calling Go's greet function with input: "${inputString}"`);
    
    // a. 将 JS 字符串编码为 UTF-8 字节
    const encoder = new TextEncoder();
    const encodedString = encoder.encode(inputString);
    const inputSize = encodedString.length;

    // b. 调用 Go 导出的 `malloc` 在 WASM 内存中分配空间
    const inputPtr = malloc(inputSize);
    if (inputPtr === 0) {
        throw new Error("malloc failed");
    }

    // c. 创建一个指向 WASM 内存的视图，并将编码后的字符串写入
    const wasmMemory = new Uint8Array(memory.buffer, inputPtr, inputSize);
    wasmMemory.set(encodedString);

    // d. 调用 Go 导出的 `greet` 函数，传递指针和大小
    const resultPtrAndSize = greet(inputPtr, inputSize);

    // --- 从 Go 接收返回的字符串 ---

    // a. 从 64 位返回值中解析出指针和大小
    // BigInt is needed because JS numbers are 53-bit precision floats.
    const resultPtr = Number(BigInt.asUintN(32, resultPtrAndSize >> 32n));
    const resultSize = Number(BigInt.asUintN(32, resultPtrAndSize));

    // b. 创建一个指向 WASM 内存中结果字符串的视图
    const resultMemory = new Uint8Array(memory.buffer, resultPtr, resultSize);
    
    // c. 将 UTF-8 字节解码为 JS 字符串
    const decoder = new TextDecoder();
    const resultString = decoder.decode(resultMemory);
    
    console.log(`Received response from Go: "${resultString}"`);

    // 5. 释放 WASM 内存
    // 释放我们为输入和输出分配的两块内存
    free(inputPtr, inputSize);
    free(resultPtr, resultSize);
}

run().catch(err => {
    console.error(err);
    process.exit(1);
}); 