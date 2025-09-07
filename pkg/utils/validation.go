package utils

import (
	"regexp"
	"strconv"
)

func IsValidObjectID(id string) bool {
	// Verifica se o ID tem 24 caracteres hexadecimais
	objectIDPatternRegex := `^[a-fA-F0-9]{24}$`
	matched, _ := regexp.MatchString(objectIDPatternRegex, id)
	return matched
}

func ValidCNPJ(cnpj string) bool {
	// Verifica se a string está vazia
	if cnpj == "" {
		return false
	}

	// Remove caracteres não numéricos
	cleanCnpj := CleanCNPJ(cnpj)

	//Verifica se tem 14 dígitos
	if len(cnpj) != 14 {
		return false
	}

	if _, err := strconv.Atoi(cleanCnpj); err != nil {
		return false
	}

	//Verifica se todos os dígitos são iguais
	if has := hasDifferentDigits(cnpj); !has {
		return false
	}

	return true
}

func CleanCNPJ(cnpj string) string {
	//Remove todos os caracteres não numéricos
	reg := regexp.MustCompile(`[^0-9]`)
	return reg.ReplaceAllString(cnpj, "")
}

func hasDifferentDigits(cnpj string) bool {
	if len(cnpj) == 0 {
		return false
	}
	firstDigit := cnpj[0]
	for i := 1; i < len(cnpj); i++ {
		if cnpj[i] != firstDigit {
			return true
		}
	}
	return false
}

func FormatCNPJ(cnpj string) string {
	clean := CleanCNPJ(cnpj)
	if len(clean) != 14 {
		return cnpj
	}
	return clean[:2] + "." + clean[2:5] + "." + clean[5:8] + "/" + clean[8:12] + "-" + clean[12:]
}
