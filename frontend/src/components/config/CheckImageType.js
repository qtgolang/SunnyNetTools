export function getImageType(data) {
    if (!(data instanceof Uint8Array)) {
        return null;
    }
    if (data.length < 20) {
        return null
    }
    if (isPNG(data)) {
        return 'png';
    }

    // 检查 JPG 格式
    if (isJPG(data)) {
        return 'jpg';
    }

    // 检查 ICO 格式
    if (isICO(data)) {
        return 'ico';
    }

    // 检查 GIF 格式
    if (isGIF(data)) {
        return 'gif';
    }

    // 检查 WEBP 格式
    if (isWEBP(data)) {
        return 'webp';
    }

    // 检查 BMP 格式
    if (isBMP(data)) {
        return 'bmp';
    }

    // 检查 TIFF 格式
    if (isTIFF(data)) {
        return 'tiff';
    }

    // 检查 SVG 格式
    if (isSVG(data)) {
        return 'svg';
    }

    // 检查 PSD 格式
    if (isPSD(data)) {
        return 'psd';
    }

    // 检查 EPS 格式
    if (isEPS(data)) {
        return 'eps';
    }

    // 检查 HEIF 格式
    if (isHEIF(data)) {
        return 'heif';
    }
    // 检查 avif 格式
    if (isAVIF(data)) {
        return 'avif';
    }

    return null;
}

function isPNG(data) {
    // PNG 图像标头是 [0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A]
    const pngHeader = [0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A];

    // 检查前 8 个字节是否匹配 PNG 标头
    for (let i = 0; i < 8; i++) {
        if (data[i] !== pngHeader[i]) {
            return false;
        }
    }

    return true;
}
function isAVIF(data) {
    // AVIF 文件的 `ftyp` 头部应该在 4~8 字节处
    return (
        data[4] === 0x66 && data[5] === 0x74 && data[6] === 0x79 && data[7] === 0x70 && // "ftyp"
        ((data[8] === 0x61 && data[9] === 0x76 && data[10] === 0x69 && data[11] === 0x66) ||  // "avif"
            (data[8] === 0x61 && data[9] === 0x76 && data[10] === 0x69 && data[11] === 0x73))    // "avis"
    );
}

function isJPG(data) {
    // JPG 图像标头是 [0xFF, 0xD8, 0xFF]
    return data[0] === 0xFF && data[1] === 0xD8 && data[2] === 0xFF;

}

function isICO(data) {
    // ICO 图像标头是 [0x00, 0x00, 0x01, 0x00]
    return data[0] === 0x00 && data[1] === 0x00 && data[2] === 0x01 && data[3] === 0x00;

}

function isGIF(data) {
    // GIF 图像标头是 [0x47, 0x49, 0x46, 0x38]
    return data[0] === 0x47 && data[1] === 0x49 && data[2] === 0x46 && data[3] === 0x38;

}

function isWEBP(data) {
    // WEBP 图像标头是 [0x52, 0x49, 0x46, 0x46, 0x??, 0x??, 0x??, 0x??, 0x57, 0x45, 0x42, 0x50]
    return data[0] === 0x52 && data[1] === 0x49 && data[2] === 0x46 && data[3] === 0x46 && data[8] === 0x57 && data[9] === 0x45 && data[10] === 0x42 && data[11] === 0x50;

}

function isBMP(data) {
    // BMP 图像标头是 [0x42, 0x4D]
    return data[0] === 0x42 && data[1] === 0x4D;

}

function isTIFF(data) {
    // TIFF 图像标头是 [0x49, 0x49, 0x2A, 0x00] 或 [0x4D, 0x4D, 0x00, 0x2A]
    return (data[0] === 0x49 && data[1] === 0x49 && data[2] === 0x2A && data[3] === 0x00) ||
        (data[0] === 0x4D && data[1] === 0x4D && data[2] === 0x00 && data[3] === 0x2A);

}

function isSVG(data) {
    // SVG 图像标头是 [0x3C, 0x3F, 0x78, 0x6D, 0x6C]
    return data[0] === 0x3C && data[1] === 0x3F && data[2] === 0x78 && data[3] === 0x6D && data[4] === 0x6C;

}

function isPSD(data) {
    // PSD 图像标头是 [0x38, 0x42, 0x50, 0x53]
    return data[0] === 0x38 && data[1] === 0x42 && data[2] === 0x50 && data[3] === 0x53;

}

function isEPS(data) {
    // EPS 图像标头是 [0x25, 0x21, 0x50, 0x53]
    return data[0] === 0x25 && data[1] === 0x21 && data[2] === 0x50 && data[3] === 0x53;

}

function isHEIF(data) {
    // HEIF 图像标头是 [0x66, 0x74, 0x79, 0x70, 0x68, 0x65, 0x69, 0x63]
    return data[0] === 0x66 && data[1] === 0x74 && data[2] === 0x79 && data[3] === 0x70 &&
        data[4] === 0x68 && data[5] === 0x65 && data[6] === 0x69 && data[7] === 0x63;

}