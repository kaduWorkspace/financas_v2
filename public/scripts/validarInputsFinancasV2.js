var valor_taxa_anual = document.getElementById('valor_taxa_anual')
var valor_taxa_anual_input = document.getElementById('valor_taxa_anual_input')
var valor_aporte = document.getElementById('valor_aporte')
var valor_aporte_input = document.getElementById('valor_aporte_input')
var valor_inicial = document.getElementById('valor_inicial')
var valor_inicial_input = document.getElementById('valor_inicial_input')
var dias_liquidez = document.getElementById('tipo_dias_liquidez_por_ano')
var form = document.getElementById('formulario_calcular')
var inputsPossiveis = [...form.elements].filter(input => !input.dataset.ignore_input)
var inputs_por_nome = inputsPossiveis.reduce((acc, curr) => {
    acc[curr.name] = curr
    return acc
}, {})
var data_final_opcoes = document.getElementById('data_final_opcao')
var data_final_especifico_wrapper = document.getElementById('data_especifica_wrapper')
var data_final_especifico_input = document.getElementById('data_final')
data_final_opcoes.addEventListener('change', (e) => {
    e.target.value == "data_especifica"
        ? data_final_especifico_wrapper.classList.remove('hidden')
        : data_final_especifico_wrapper.classList.add('hidden')
})
var init = () => {
    if(!form) return
    var inputs = form.querySelectorAll('input');
    form.addEventListener("input", e => {
        window.handleErrorsEvent(e);
        window.processarInputs();
        // Store input values in sessionStorage
        inputs.forEach(input => {
            sessionStorage.setItem(input.name || input.id, input.value);
        });
    });
    // Restore input values from sessionStorage on page load
    inputs.forEach(input => {
        var storedValue = sessionStorage.getItem(input.name || input.id);
        if (storedValue) {
            input.value = storedValue;
        }
    });
    window.processarInputs();
    if(data_final_opcoes.value !== "data_especifica") {
        var tipo = data_final_opcoes.value == "6" ? "meses" : "anos";
        var data_resultado = window.incrementar_data(parseInt(data_final_opcoes.value), tipo)
        data_final_especifico_input.value = data_resultado
    }
    var data_inicial_input = document.getElementById("data_inicial")
    data_inicial_input.value = window.formatar_data(new Date())
    for (var input of form.elements) {
        if (input.id) {
            var savedValue = sessionStorage.getItem(input.id);
            if (savedValue) {
                input.value = savedValue;
            }
        }
    }
    if(data_final_opcoes.value == "data_especifica") {
        data_final_especifico_wrapper.classList.remove('hidden')
    }

    valor_aporte.addEventListener('input', window.evento_mascara_monetaria)
    valor_inicial.addEventListener('input', window.evento_mascara_monetaria)
    valor_taxa_anual.addEventListener('input', e => {
        e.target.value = window.mascara_monetaria(e.target.value)
    })
    document.body.addEventListener("htmx:configRequest", event => {
        if(event.detail.elt.id === "formulario_calcular") {
            if (!window.validarRequest(event)) {
                event.preventDefault();
            }
        }
    });
}
init()

