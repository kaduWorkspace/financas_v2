const tipoAumentoFrequenciaInput = document.getElementById('tipo_aumento_frequencia')
const valorAumentoAporte = document.getElementById('valor_aumento_aporte')
const valorAumentoAporteWrapper = document.getElementById('valor_aumento_aporte_wrapper')
const form = document.getElementById('formulario_calcular')
const inputsPossiveis = [...form.elements].filter(input => !input.dataset.ignore_input)
const data_final_opcoes = document.getElementById('data_final_opcao')
const data_final_especifico_wrapper = document.getElementById('data_especifica_wrapper')
const data_final_especifico_input = document.getElementById('data_final')
data_final_opcoes.addEventListener('change', ({target:{value}}) => {
    value == "data_especifica"
        ? data_final_especifico_wrapper.classList.remove('hidden')
        : data_final_especifico_wrapper.classList.add('hidden')
})
const formatar_data = (data) => {
    const ano = data.getFullYear();
    const mes = String(data.getMonth() + 1).padStart(2, "0"); // Mês começa em 0
    const dia = String(data.getDate()).padStart(2, "0");
    const data_resultado = `${ano}-${mes}-${dia}`;
    return data_resultado
}
const incrementar_data = (quantidade = 6, tipo = "meses") => {
    const data = new Date();
    tipo == "meses"
        ? data.setMonth(data.getMonth() + quantidade)
        : data.setFullYear(data.getFullYear() + quantidade)
    return formatar_data(data);
}
const buscarValoresInput = () => {
    return inputsPossiveis.reduce((acc, item) => {
            acc[item.name] = item.type == 'number' ? parseInt(item.value) : item.value
            return acc;
    }, {})
}
const validarValoresInputs = (inputs, validarNull = false) => {
    const {valor_inicial, aporte_mensal, aporte_semestral, data_final, tipo_frequencia_aporte, valor_aumento_aporte} = inputs
    const erros = [];
    if(valor_inicial < 0 || (validarNull && ["",null,false].includes(valor_inicial))) erros.push(["valor_inicial","Valor inicial inválido"]);
    if(aporte_mensal < 0 || (validarNull && ["",null,false].includes(aporte_mensal))) erros.push( ["aporte_mensal","Aporte mensal inválido"]);
    if(aporte_semestral < 0 || (validarNull && ["",null,false].includes(aporte_semestral))) erros.push( ["aporte_semestral","Aporte semestral inválido"]);
    if(valor_aumento_aporte < 0 || (validarNull && ["",null,false].includes(valor_aumento_aporte))) erros.push( ["valor_aumento_aporte","Aumento inválido"]);
    if(aporte_mensal > 1000000) erros.push(["aporte_mensal","Aporte mensal muito alto"]);
    if(aporte_semestral > 1000000) erros.push(["aporte_semestral","Aporte semestral muito alto"]);
    if(valor_aumento_aporte > 100000) erros.push(["valor_aumento_aporte","Aumento muito alto"]);
    if(!data_final) erros.push(["data_final","Data final inválida"]);
    return erros.length ? erros : false;
}
tipoAumentoFrequenciaInput.addEventListener('input', event => {
    if(tipoAumentoFrequenciaInput.value == "false") {
        valorAumentoAporteWrapper.classList.add('hidden')
        valorAumentoAporte.value = ""
    }
    else valorAumentoAporteWrapper.classList.remove('hidden')
})
inputsPossiveis.forEach(item => {
    item.addEventListener('input', event => {
        const errorSpan = document.getElementById(`error_${item.id}`)
        if(errorSpan)
            errorSpan.classList.add('hidden');
        let validacao = validarValoresInputs(buscarValoresInput(), false);
        if(!validacao) return;
        const validacaoInput = validacao.find(valid => valid[0] == item.name);
        if(!validacaoInput) return;
        errorSpan.innerText = validacaoInput[1];
        errorSpan.classList.remove('hidden')
        return;
    });
})
document.addEventListener("DOMContentLoaded", () => {
    if(!form) return
    const data_inicial_input = document.getElementById("data_inicial")
    data_inicial_input.value = formatar_data(new Date())
    for (const input of form.elements) {
        if (input.name) {
            const savedValue = sessionStorage.getItem(input.name);
            if (savedValue) {
                input.value = savedValue;
            }
        }
    }
    if(data_final_opcoes.value == "data_especifica") {
        data_final_especifico_wrapper.classList.remove('hidden')
    }
    if(tipoAumentoFrequenciaInput.value != "false") {
        valorAumentoAporteWrapper.classList.remove('hidden')
    }
    form.addEventListener("input", function (event) {
        const { name, value } = event.target;
        if (name) {
            sessionStorage.setItem(name, value);
        }
    });
    form.addEventListener('submit', event => {
        event.preventDefault();
        if(data_final_opcoes.value !== "data_especifica") {
            const tipo = data_final_opcoes.value == "6" ? "meses" : "anos";
            data_final_especifico_input.value = incrementar_data(parseInt(data_final_opcoes.value), tipo)
        }
        const validacoes = validarValoresInputs(buscarValoresInput(), true);
        if(validacoes) {
            validacoes.forEach(validacao => {
                const errorSpan = document.getElementById(`error_${validacao[0]}`)
                errorSpan.innerText = validacao[1];
                errorSpan.classList.remove('hidden')
            })
            return
        }
        inputsPossiveis.filter(input => input.type == "number").forEach(input => {
            if(input.value === "") {
                input.value = 0.0
            }
        })
        form.submit()
    })
})
