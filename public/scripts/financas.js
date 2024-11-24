const wrapperMain = document.getElementById("wrapperGraficos")
function formatarData(data) {
    // Verifica se a entrada é um objeto Date, caso contrário tenta convertê-la
    if (!(data instanceof Date)) {
        data = new Date(data);
    }

    // Formata a data no padrão DD/MM/YYYY
    const dia = String(data.getDate()).padStart(2, '0'); // Adiciona zero à esquerda se necessário
    const mes = String(data.getMonth() + 1).padStart(2, '0'); // Meses começam do zero
    const ano = data.getFullYear();

    return `${dia}/${mes}/${ano}`;
}
const tipoAumentoFrequenciaInput = document.getElementById('tipo-aumento-frequencia')
const valorAumentoAporte = document.getElementById('valor_aumento_aporte')
tipoAumentoFrequenciaInput.addEventListener('input', event => {
    if(tipoAumentoFrequenciaInput.value == "false") {
        valorAumentoAporte.classList.add('hidden')
        valorAumentoAporte.value = ""
    }
    else valorAumentoAporte.classList.remove('hidden')
})
let dados = [];
const calcular = async ({valor_inicial, aporte_mensal, aporte_semestral, data_final, tipo_frequencia_aporte, valor_aumento_aporte}) => {
    let dataIni = new Date()
    const data_inicial = `${dataIni.getDate()}/${dataIni.getMonth()}/${dataIni.getFullYear()}`
    const body = {valor_inicial, aporte_mensal, aporte_semestral, data_final: formatarData(data_final), tipo_frequencia_aporte, valor_aumento_aporte, data_inicial}
    const request = await fetch('/calcular-juros', {
        method: "POST",
        headers: {'content-type':'application/json'},
        body: JSON.stringify(body)
    })
    console.log(request.status)
    const res = await request.json();
    return request.status == 200 ? res : null;
}
const buscarValoresInput = () => {
    return [
        ...document.getElementsByTagName('input'),
        ...document.getElementsByTagName('select') ]
        .reduce((acc, item) => {
            acc[item.name] = item.type == 'number' ? parseInt(item.value) : item.value
            return acc;
    }, {})
}
document.getElementById('calcular').addEventListener('click', async event => {
    dados = await calcular(buscarValoresInput())
    if(!dados) {
        dados = [];
        return;
    }
    if(wrapperMain.classList.contains('hidden')) wrapperMain.classList.toggle('hidden')
    console.log({dados})
})
const formatarValorMonetario = (valor) => `R$ ${(Math.floor(valor*100)/100).toLocaleString('pt-br',{style: 'currency', currency: 'BRL'})}`;

const divGraficos = document.getElementById('grafico')

const cabecalhoTabela = titulo => `
        <th scope="col" class="px-6 py-3">
            ${titulo}
        </th>
`;
const montaCabecalhoTabela = titulos => {
    const cabecalhos = titulos.map(cabecalhoTabela).join('')
    return `<tr>${cabecalhos}</tr>`
}
const montaItemLinha = dado => `<th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
${dado}
    </th>`

const montaLinhaTabela = (dados) => `
        <tr class="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
${dados.map(montaItemLinha).join('')}
        </tr>`

document.getElementById('dias-grafico').addEventListener('click', event => {
    if(divGraficos.classList.contains('hidden')) divGraficos.classList.toggle('hidden')
    const cabecalho = montaCabecalhoTabela(["Data", "Valorização", "Resultado"])
    const dados_ = dados.dias.reduce((acc, curr) => {
        const {data, valorizacao, resultado_com_valorizacao} = curr
        acc.push([data, formatarValorMonetario(valorizacao), formatarValorMonetario(resultado_com_valorizacao)])
        return acc
    }, [])
    const linhasTabela = dados_.map(montaLinhaTabela).join('')
    const tabela = `
        <div class="relative overflow-x-auto">
            <table class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    ${cabecalho}
                </thead>
                <tbody>
                    ${linhasTabela}
                </tbody>
            </table>
        </div>
    `
    document.getElementById('resultadoGrafico').innerHTML = tabela
})
document.getElementById('meses-grafico').addEventListener('click', event => {
    if(divGraficos.classList.contains('hidden')) divGraficos.classList.toggle('hidden')
    const cabecalho = montaCabecalhoTabela(["Data", "Valorização", "Resultado"])
    const dados_ = dados.meses.reduce((acc, curr) => {
        const {data_inicial, data_final, valorizacao, resultado_com_valorizacao} = curr
        acc.push([`${data_inicial} - ${data_final}`, formatarValorMonetario(valorizacao), formatarValorMonetario(resultado_com_valorizacao)])
        return acc
    }, [])
    const linhasTabela = dados_.map(montaLinhaTabela).join('')
    const tabela = `
        <div class="relative overflow-x-auto">
            <table class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    ${cabecalho}
                </thead>
                <tbody>
                    ${linhasTabela}
                </tbody>
            </table>
        </div>
    `
    document.getElementById('resultadoGrafico').innerHTML = tabela
})
document.getElementById('semestres-grafico').addEventListener('click', event => {
    if(divGraficos.classList.contains('hidden')) divGraficos.classList.toggle('hidden')
    const cabecalho = montaCabecalhoTabela(["Data", "gasto", "Valorização", "Resultado"])
    const dados_ = dados.semestres.reduce((acc, curr) => {
        const {data_inicial, data_final, valorizacao, resultado_com_valorizacao, gasto} = curr
        acc.push([`${data_inicial} - ${data_final}`, formatarValorMonetario(gasto),formatarValorMonetario(valorizacao), formatarValorMonetario(resultado_com_valorizacao)])
        return acc
    }, [])
    const linhasTabela = dados_.map(montaLinhaTabela).join('')
    const tabela = `
        <div class="relative overflow-x-auto">
            <table class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    ${cabecalho}
                </thead>
                <tbody>
                    ${linhasTabela}
                </tbody>
            </table>
        </div>
    `
    document.getElementById('resultadoGrafico').innerHTML = tabela
})

document.getElementById('anos-grafico').addEventListener('click', event => {
    if(divGraficos.classList.contains('hidden')) divGraficos.classList.toggle('hidden')
    const cabecalho = montaCabecalhoTabela(["Data", "gasto", "Valorização", "Resultado"])
    const dados_ = dados.anos.reduce((acc, curr) => {
        const {data_inicial, valorizacao, resultado_com_valorizacao, gasto} = curr
        acc.push([data_inicial, formatarValorMonetario(gasto), formatarValorMonetario(valorizacao), formatarValorMonetario(resultado_com_valorizacao)])
        return acc
    }, [])
    const linhasTabela = dados_.map(montaLinhaTabela).join('')
    const tabela = `
        <div class="relative overflow-x-auto">
            <table class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    ${cabecalho}
                </thead>
                <tbody>
                    ${linhasTabela}
                </tbody>
            </table>
        </div>
    `
    document.getElementById('resultadoGrafico').innerHTML = tabela
})
document.getElementById('fecharGrafico').addEventListener('click', event => {
    if(!divGraficos.classList.contains('hidden')) divGraficos.classList.toggle('hidden')
})
