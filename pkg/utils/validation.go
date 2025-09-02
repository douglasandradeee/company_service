package utils

import (
	"regexp"
	"strconv"
)

func ValidCNPJ(cnpj string) bool {
	//Remove caracteres não numéricos
	cnpj = CleanCNPJ(cnpj)

	//Verifica se tem 14 dígitos
	if len(cnpj) != 14 {
		return false
	}

	//Verifica se todos os dígitos são iguais
	if allDigitsEqual(cnpj) {
		return false
	}

	//Verifica se os dígitos verificadores são válidos
	d1 := calculateCNPJCheckerDigit(cnpj, 1)
	d2 := calculateCNPJCheckerDigit(cnpj, 2)
	if cnpj[12:14] != d1+d2 {
		return false
	}

	return true
}

func CleanCNPJ(cnpj string) string {
	//Remove todos os caracteres não numéricos
	reg := regexp.MustCompile(`[^0-9]`)
	return reg.ReplaceAllString(cnpj, "")
}

func allDigitsEqual(cnpj string) bool {
	firstDigit := cnpj[0]
	for i := 1; i < len(cnpj); i++ {
		if cnpj[i] != firstDigit {
			return false
		}
	}
	return true
}

func calculateCNPJCheckerDigit(cnpj string, digitNumber int) string {
	weights := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0

	// Para o primeiro dígito verificador: usa os 12 primeiros dígitos
	// Para o segundo dígito verificador: usa os 13 primeiros dígitos
	numDigits := 12
	if digitNumber == 2 {
		numDigits = 13
	}

	// Garantir que não ultrapassamos o tamanho do CNPJ
	if numDigits > len(cnpj) {
		numDigits = len(cnpj)
	}

	for i := 0; i < numDigits; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))

		// Ajustar o índice do peso baseado no dígito verificador
		weightIndex := i
		if digitNumber == 2 {
			// Para o segundo dígito, os pesos começam do índice 1
			weightIndex = i + 1
		}

		// Garantir que o weightIndex está dentro dos limites
		if weightIndex < len(weights) {
			sum += digit * weights[weightIndex]
		}
	}

	remainder := sum % 11
	if remainder < 2 {
		return "0"
	}
	return strconv.Itoa(11 - remainder)
}

func FormatCNPJ(cnpj string) string {
	clean := CleanCNPJ(cnpj)
	if len(clean) != 14 {
		return cnpj
	}
	return clean[:2] + "." + clean[2:5] + "." + clean[5:8] + "/" + clean[8:12] + "-" + clean[12:]
}
