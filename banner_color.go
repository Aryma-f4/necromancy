package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Konversi hex color ke ANSI color
type ANSIColor struct {
	Code  string
	Reset string
}

var Reset = "\033[0m"

func hexToANSIColor(hexColor string) string {
	hexColor = strings.TrimPrefix(hexColor, "#")
	if len(hexColor) != 6 {
		return ""
	}

	// Parse RGB
	r, err1 := strconv.ParseInt(hexColor[0:2], 16, 0)
	g, err2 := strconv.ParseInt(hexColor[2:4], 16, 0)
	b, err3 := strconv.ParseInt(hexColor[4:6], 16, 0)

	if err1 != nil || err2 != nil || err3 != nil {
		return ""
	}

	// Sederhanakan ke warna dasar ANSI berdasarkan dominasi warna
	ri, gi, bi := int(r), int(g), int(b)

	// Hitung kecerahan
	brightness := (ri + gi + bi) / 3

	// Tentukan warna dominan
	if ri > gi && ri > bi {
		// Dominan merah
		if brightness > 180 {
			return "\033[91m" // Bright red
		}
		return "\033[31m" // Red
	} else if gi > ri && gi > bi {
		// Dominan hijau
		if brightness > 180 {
			return "\033[92m" // Bright green
		}
		return "\033[32m" // Green
	} else if bi > ri && bi > gi {
		// Dominan biru
		if brightness > 180 {
			return "\033[94m" // Bright blue
		}
		return "\033[34m" // Blue
	} else {
		// Netral atau grayscale
		if brightness > 230 {
			return "\033[97m" // Bright white
		} else if brightness > 180 {
			return "\033[37m" // White
		} else if brightness > 120 {
			return "\033[90m" // Gray
		} else {
			return "\033[30m" // Black
		}
	}
}

func parseBBCodeColor(text string) string {
	// Pattern untuk menangkap tag [color=#hex]content[/color]
	colorPattern := regexp.MustCompile(`\[color=([^\]]+)\](.*?)\[/color\]`)

	result := text
	matches := colorPattern.FindAllStringSubmatch(result, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			colorCode := match[1]
			content := match[2]

			ansiColor := hexToANSIColor(colorCode)
			if ansiColor != "" {
				coloredContent := ansiColor + content + Reset
				result = strings.Replace(result, match[0], coloredContent, 1)
			} else {
				// Jika warna tidak bisa dikonversi, hapus tagnya saja
				result = strings.Replace(result, match[0], content, 1)
			}
		}
	}

	// Hapus tag BBCode yang tidak diproses
	result = regexp.MustCompile(`\[color=[^\]]+\]`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[/color\]`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[size=[^\]]+\]`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[/size\]`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[font=[^\]]+\]`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`\[/font\]`).ReplaceAllString(result, "")

	// Hapus tag HTML span jika ada
	result = regexp.MustCompile(`<span[^>]*>`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`</span>`).ReplaceAllString(result, "")

	return result
}

func printFallbackBanner() {
	// Banner fallback sederhana
	fmt.Println("╔═══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        NECROMANCY - Advanced Shell                    ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════════════╝")
}

func printColoredBanner() {
	// Baca file ascii.txt
	file, err := os.Open("ascii.txt")
	if err != nil {
		// Jika file tidak ditemukan, gunakan banner fallback
		printFallbackBanner()
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse warna dari BBCode
		coloredLine := parseBBCodeColor(line)
		lines = append(lines, coloredLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading banner file:", err)
		printFallbackBanner()
		return
	}

	// Tampilkan banner
	for _, line := range lines {
		fmt.Println(line)
	}

	// Tampilkan versi dan link
	fmt.Println("")
	fmt.Println("\033[90m v1.5.0 Stable Release - Advanced Post Exploitation Tools \033[0m")
	fmt.Println("\033[90m https://github.com/Aryma-f4/necromancy \033[0m")
	fmt.Println("")
}
