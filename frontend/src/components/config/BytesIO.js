class BytesIO {
    /**
     * 构造函数，初始化数据缓冲区
     * @param {Uint8Array} data 输入的二进制数据
     */
    constructor(data) {
        this.array = new Uint8Array(data); // 存储数据
        this.pos = 0; // 当前读取位置
    }
    /**
     * 读取直到分隔符的数据，跳过分隔符可选
     * @param {number} delim - 分隔符的 ASCII 码
     * @param {boolean} excludeDelim - 是否排除分隔符（如果为 true 则不包含分隔符）
     * @returns {Uint8Array|null} 读取到的数据
     */
    ReadBytes(delim, excludeDelim) {
        if (this.pos >= this.array.length) return null; // EOF

        let start = this.pos;
        while (this.pos < this.array.length) {
            if (this.array[this.pos] === delim) {
                let result = this.array.slice(start, this.pos); // 不包括分隔符
                if (!excludeDelim) {
                    result = this.array.slice(start, this.pos + 1); // 如果不排除分隔符，包含分隔符
                }
                this.pos++; // 跳过分隔符
                return result;
            }
            this.pos++;
        }
        return this.array.slice(start, this.pos); // 没找到 delim，返回剩余数据
    }

    /**
     * 读取一行数据，跳过换行符和回车符可选
     * @param {boolean} excludeNewline - 是否排除 `\r` 和 `\n`（如果为 true 则不包含换行符）
     * @returns {Uint8Array|null} 读取到的行数据
     */
    ReadLine(excludeNewline) {
        if (this.pos >= this.array.length) return null; // EOF

        let start = this.pos;
        while (this.pos < this.array.length) {
            if (this.array[this.pos] === 10) { // `\n` 的 ASCII 码
                let line = this.array.slice(start, this.pos); // 不包括 `\n`
                if (excludeNewline && line.length > 0 && line[line.length - 1] === 13) { // 检查是否是 `\r`
                    line = line.slice(0, -1); // 去掉 `\r`
                }
                this.pos++; // 跳过 `\n`
                return line;
            }
            this.pos++;
        }
        return this.array.slice(start, this.pos); // 返回最后一行（可能没有换行符）
    }

    /**
     * 读取剩余的所有数据
     * @returns {Uint8Array} 剩余的数据
     */
    ReadAll() {
        let remaining = this.array.slice(this.pos);
        this.pos = this.array.length; // 读完了
        return remaining;
    }

    /**
     * 获取当前缓冲区中剩余的可读字节数
     * @returns {number} 剩余字节数
     */
    Buffered() {
        return this.array.length - this.pos;
    }

    /**
     * 读取 `p.length` 个字节到 `p`，返回实际读取的字节数
     * @param {Uint8Array} p 目标缓冲区
     * @returns {number} 实际读取的字节数
     */
    Read(p) {
        let n = Math.min(p.length, this.Buffered());
        if (n === 0) return 0; // 没有数据可读
        p.set(this.array.slice(this.pos, this.pos + n));
        this.pos += n;
        return n;
    }

    /**
     * 预览接下来的 `n` 个字节，不改变 `pos`
     * @param {number} n 需要预览的字节数
     * @returns {Uint8Array|null} 预览的数据，如果超出范围返回 `null`
     */
    Peek(n) {
        if (n > this.Buffered()) return null; // 超过范围
        return this.array.slice(this.pos, this.pos + n);
    }

    /**
     * 丢弃 `n` 个字节，相当于跳过 `n` 个字节
     * @param {number} n 需要跳过的字节数
     * @returns {number} 丢弃后的当前位置
     */
    Discard(n) {
        this.pos = Math.min(this.pos + n, this.array.length);
        return this.pos;
    }

    /**
     * 撤回最近一次读取的字节（相当于回退 1 字节）
     * 如果已经回退到起始位置，则抛出错误
     * @throws {Error} 如果已经到达缓冲区起始位置
     */
    UnreadByte() {
        if (this.pos > 0) {
            this.pos--;
        } else {
            this.pos = 0;
        }
    }

    /**
     * 撤回最近一次读取的 UTF-8 字符（支持多字节字符）
     * 由于 UTF-8 是变长编码，需要检查前缀字节
     * @throws {Error} 如果已经到达缓冲区起始位置
     */
    UnreadRune() {
        if (this.pos > 0) {
            do {
                this.pos--;
            } while (this.pos > 0 && (this.array[this.pos] & 0xC0) === 0x80);
        } else {
            this.pos = 0;
        }
    }
}

export function NewBytesIO(data) {
    return new BytesIO(data)
}