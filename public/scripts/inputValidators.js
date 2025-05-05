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

const validarValoresInputs = (validarNull = false) => {
    const erros = [];
    let valor_inicial_v = tirar_mascara(mascara_monetaria(tirar_mascara(valor_inicial.value)))
    let valor_aporte_v = tirar_mascara(mascara_monetaria(tirar_mascara(valor_aporte.value)))
    if(valor_inicial_v < 0 || (validarNull && ["",null,false].includes(valor_inicial))) erros.push(["error_valor_inicial","Valor inicial inválido"]);
    if(valor_aporte_v > 1000000000) erros.push(["error_valor_aporte","Aporte mensal muito alto"]);
    if(!(valor_aporte_v + valor_inicial_v > 0)) erros.push(
        ["error_valor_inicial","O valor inicial ou valor de aporte devem ser preenchidos!"],
        ["error_valor_aporte","O valor inicial ou valor de aporte devem ser preenchidos!"]
    );
    if(!document.getElementById("data_final").value) erros.push(["error_data_final","Data final inválida"]);
    return erros.length ? erros : false;
}
const processarInputs = () => {
    inputsPossiveis.filter(input => input.type == "number").forEach(input => {
        if(input.value === "") {
            input.value = 0.0
            console.log(`Empty input detected for ${input.name}, defaulting to 0.0`)
        }
    })

    let taxa_anual_v = Number(valor_taxa_anual.value.replaceAll('.','').replaceAll(',','.').replaceAll('%','').trim()) ?? 0;
    let valor_aporte_v = Number(valor_aporte.value.replaceAll('.','').replaceAll(',','.').replaceAll('R$','').trim()) ?? 0;
    let valor_inicial_v = Number(valor_inicial.value.replaceAll('.','').replaceAll(',','.').replaceAll('R$','').trim()) ?? 0;

    if(data_final_opcoes.value !== "data_especifica") {
        const tipo = data_final_opcoes.value == "6" ? "meses" : "anos";
        const data_resultado = incrementar_data(parseInt(data_final_opcoes.value), tipo)
        data_final_especifico_input.value = data_resultado
    }
    valor_aporte_input.value = !!valor_aporte_v ? valor_aporte_v : 0;
    valor_taxa_anual_input.value = !!taxa_anual_v ? taxa_anual_v : 0;
    valor_inicial_input.value = !!valor_inicial_v ? valor_inicial_v : 0;
    return;
}
const handleErrorsEvent = e => {
    const errorSpan = document.getElementById(`error_${e.target.id}`)
    if(errorSpan)
        errorSpan.classList.add('hidden');
    const validacao = validarValoresInputs(false);
    if(!validacao) {
        return;
    }
    const validacaoInput = validacao.find(([error_span_target_name]) => error_span_target_name == errorSpan.id);
    if(!validacaoInput) return;
    errorSpan.innerText = validacaoInput[1];
    errorSpan.classList.remove('hidden')
    return;
}
const validarRequest = event => {
    //processarInputs();
    const validacoes = validarValoresInputs(true);
    if(validacoes) {
        console.log(validacoes)
        validacoes.forEach(validacao => {
            const errorSpan = document.getElementById(validacao[0])
            errorSpan.innerText = validacao[1];
            errorSpan.classList.remove('hidden')
        })
        console.log("retornando false")
        return false;
    }
    return true;
}
const tirar_mascara = (v, p) => {
    return v.replace(/\D/g, '') ?? 0
}
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
    e.target.value = mascara_monetaria(e.target.value)
}

// Assign all functions to window object
window.formatar_data = formatar_data;
window.incrementar_data = incrementar_data;
window.validarValoresInputs = validarValoresInputs;
window.processarInputs = processarInputs;
window.handleErrorsEvent = handleErrorsEvent;
window.validarRequest = validarRequest;
window.tirar_mascara = tirar_mascara;
window.to_valor_monetario = to_valor_monetario;
window.removerZerosIniciais = removerZerosIniciais;
window.mascara_monetaria = mascara_monetaria;
window.evento_mascara_monetaria = evento_mascara_monetaria;
