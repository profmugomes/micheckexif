// Copyright (C) 2024-2026 Murilo Gomes <profmugomes.com.br>
// SPDX-License-Identifier: GPL-2.0-only

// Site: https://www.profmugomes.com.br

package controls

import (
	"github.com/profmugomes/mglang/v3"
)

func LoadTranslations() {
	lang := mglang.GetLang()
	switch lang {
	case "pt":
		mglang.Set("About", "Sobre")
		mglang.Set("Support MiCheckExif", "Apoie MiCheckExif")
		mglang.Set("About MiCheckExif", "Sobre MiCheckExif")
		mglang.Set("Files", "Arquivos")
		mglang.Set("File", "Arquivo")
		mglang.Set("Information", "Informação")
		mglang.Set("Add File", "Adicionar Arquivo")
		mglang.Set("Select the file(s)", "Selecione o(s) arquivo(s)")
		mglang.Set("Remove", "Remover")
		mglang.Set("Remove All", "Remover Todos")
		mglang.Set("Remove original file", "Remover arquivo original")
		mglang.Set("Yes", "Sim")
		mglang.Set("No", "Não")
		mglang.Set("Get Data", "Obter Dados")
		mglang.Set("Getting Data: %s", "Obtendo Dados: %s")
		mglang.Set("Remove Data", "Remover Dados")
		mglang.Set("Removing Data: %s", "Removendo Dados: %s")
		mglang.Set("File not changed", "Arquivo não alterado")
		mglang.Set("Data removed", "Dados removidos")
	}
}

func T(key string, args ...any) string {
	return mglang.T(key, args...)
}
