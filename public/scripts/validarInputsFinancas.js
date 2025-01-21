const tipoAumentoFrequenciaInput = document.getElementById('tipo_aumento_frequencia')
const valorAumentoAporte = document.getElementById('valor_aumento_aporte')
const valorAumentoAporteWrapper = document.getElementById('valor_aumento_aporte_wrapper')
const form = document.getElementById('formulario_calcular')
const inputsPossiveis = [...document.getElementsByTagName('input'),...document.getElementsByTagName('select')]
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
form.addEventListener('submit', event => {
    event.preventDefault();
    const validacoes = validarValoresInputs(buscarValoresInput(), true);
    if(validacoes) {
        validacoes.forEach(validacao => {
            const errorSpan = document.getElementById(`error_${validacao[0]}`)
            errorSpan.innerText = validacao[1];
            errorSpan.classList.remove('hidden')
        })
        return
    }
    form.submit()
})
