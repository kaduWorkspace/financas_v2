const valor_taxa_anual = document.getElementById('valor_taxa_anual')
const valor_taxa_anual_input = document.getElementById('valor_taxa_anual_input')
const valor_aporte = document.getElementById('valor_aporte')
const valor_aporte_input = document.getElementById('valor_aporte_input')
const valor_inicial = document.getElementById('valor_inicial')
const valor_inicial_input = document.getElementById('valor_inicial_input')
const dias_liquidez = document.getElementById('tipo_dias_liquidez_por_ano')
const form = document.getElementById('formulario_calcular')
const inputsPossiveis = [...form.elements].filter(input => !input.dataset.ignore_input)
const inputs_por_nome = inputsPossiveis.reduce((acc, curr) => {
    acc[curr.name] = curr
    return acc
}, {})
const tirar_mascara = (v, p) => {
    return Number(v.replaceAll('.','').replaceAll(',','.').replaceAll(p,'').trim()) ?? 0
}
const inputs_por_nome_valor = () => inputsPossiveis.reduce((acc, curr) => {
    acc[curr.name] = curr.value
    return acc
}, {})
const data_final_opcoes = document.getElementById('data_final_opcao')
const data_final_especifico_wrapper = document.getElementById('data_especifica_wrapper')
const data_final_especifico_input = document.getElementById('data_final')
const to_valor_monetario = string => {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(string);
}
function removerZerosIniciais(str) {
    const regex = /^0*(\d+)$/;
    const match = str.match(regex);
    if (match) {
        return match[1];
    } else {
        return str
    }
}
const mascara_monetaria = valor => {
    valor = valor.replace(/\D/g, '')
    valor = removerZerosIniciais(valor)
    let valor_split = valor.split('')
    if (valor_split.length == 0) {
        return "0,00"
    }
    if (valor_split.length == 1) {
        return `0,0${valor_split.pop()}`
    }
    if (valor_split.length == 2) {
        return `0,${valor_split.shift()}${valor_split.shift()}`
    }
    if (valor_split.length == 3) {
        return `${valor_split.shift()},${valor_split.shift()}${valor_split.shift()}`
    }
    let centavos = [valor_split.pop(), valor_split.pop()].reverse().join("");
    let grupos_de_tres = [[]];
    let cont = 0;
    let cont_grupo = 0;
    let resto = valor_split.length % 3;
    let primeiros = []
    while (resto > 0) {
        primeiros.push(valor_split.shift())
        resto--
    }
    while(cont < valor_split.length) {
        const atual = valor_split[cont]
        if(grupos_de_tres[cont_grupo].length == 3) {
            grupos_de_tres.push([])
            cont_grupo++;
        }
        grupos_de_tres[cont_grupo].push(atual);
        cont++;
    }
    let valor_formatado = grupos_de_tres.map(grupo => grupo.join('')).join('.')
    let resultado = "";
    if(primeiros.length) {
        if (valor_formatado !== "") {
            resultado = primeiros.join("") + "." + valor_formatado + "," + centavos;
        } else {
            resultado = primeiros.join("") + "," + centavos;
        }
    } else {
        resultado = valor_formatado + "," + centavos;
    }
    return resultado;
}
const evento_mascara_monetaria = e => {
    e.target.value = "R$ " + mascara_monetaria(e.target.value)
}
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
    const { valor_inicial, valor_aporte, data_final } = inputs
    const erros = [];
    if(tirar_mascara(valor_inicial, "R$") < 0 || (validarNull && ["",null,false].includes(valor_inicial))) erros.push(["valor_inicial","Valor inicial inválido"]);
    if(tirar_mascara(valor_aporte, "R$") > 1000000) erros.push(["aporte_mensal","Aporte mensal muito alto"]);
    if(!data_final) erros.push(["data_final","Data final inválida"]);
    return erros.length ? erros : false;
}
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
    form.addEventListener("input", function (event) {
        const { name, value } = event.target;
        if (name) {
            sessionStorage.setItem(name, value);
        }
    });
    valor_aporte.addEventListener('input', evento_mascara_monetaria)
    valor_inicial.addEventListener('input', evento_mascara_monetaria)
    valor_taxa_anual.addEventListener('input', e => {
        e.target.value = "% " + mascara_monetaria(e.target.value)
    })
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
        let taxa_anual_v = Number(valor_taxa_anual.value.replaceAll('.','').replaceAll(',','.').replaceAll('%','').trim()) ?? 0
        let valor_aporte_v = Number(valor_aporte.value.replaceAll('.','').replaceAll(',','.').replaceAll('R$','').trim()) ?? 0
        let valor_inicial_v = Number(valor_inicial.value.replaceAll('.','').replaceAll(',','.').replaceAll('R$','').trim()) ?? 0
        valor_aporte_input.value = !!valor_aporte_v ? valor_aporte_v : 0
        valor_taxa_anual_input.value = !!taxa_anual_v ? taxa_anual_v : 0
        valor_inicial_input.value = !!valor_inicial_v ? valor_inicial_v : 0
        form.submit()
    })
})
