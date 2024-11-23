type HSL = { h: number; s: number; l: number };

function hexToHsl(hex: string): HSL {
    // Remove # if present
    const normalizedHex = hex.replace("#", "");

    const r = parseInt(normalizedHex.substring(0, 2), 16) / 255;
    const g = parseInt(normalizedHex.substring(2, 4), 16) / 255;
    const b = parseInt(normalizedHex.substring(4, 6), 16) / 255;

    const max = Math.max(r, g, b);
    const min = Math.min(r, g, b);
    const delta = max - min;

    const l = (max + min) / 2;
    let h = 0;
    let s = 0;

    if (delta !== 0) {
        s = l > 0.5 ? delta / (2 - max - min) : delta / (max + min);

        switch (max) {
            case r:
                h = (g - b) / delta + (g < b ? 6 : 0);
                break;
            case g:
                h = (b - r) / delta + 2;
                break;
            case b:
                h = (r - g) / delta + 4;
                break;
        }

        h *= 60;
    }

    return { h, s: s * 100, l: l * 100 };
}

function hslToHex(h: number, s: number, l: number): string {
    const normalizedS = s / 100;
    const normalizedL = l / 100;

    const c = (1 - Math.abs(2 * normalizedL - 1)) * normalizedS;
    const x = c * (1 - Math.abs((h / 60) % 2 - 1));
    const m = normalizedL - c / 2;

    let r = 0,
        g = 0,
        b = 0;

    if (h < 60) {
        r = c;
        g = x;
    } else if (h < 120) {
        r = x;
        g = c;
    } else if (h < 180) {
        g = c;
        b = x;
    } else if (h < 240) {
        g = x;
        b = c;
    } else if (h < 300) {
        r = x;
        b = c;
    } else {
        r = c;
        b = x;
    }

    const toHex = (n: number) =>
        Math.round((n + m) * 255)
            .toString(16)
            .padStart(2, "0");

    return `#${toHex(r)}${toHex(g)}${toHex(b)}`;
}


export function adjustHexColor(hex: string, saturationIncrease: number, lightnessDecrease: number): string {
    const { h, s, l } = hexToHsl(hex);

    // Adjust saturation and lightness
    const newS = Math.min(100, s + saturationIncrease);
    const newL = Math.max(0, l - lightnessDecrease);

    return hslToHex(h, newS, newL);
}