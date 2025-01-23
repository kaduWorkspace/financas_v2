function formatarDatasParaTemplate(startDate, endDate) {
    const formatDate = (date) => {
        const [day, month, year] = date.split('/');
        const formattedMonth = month.padStart(2, '0'); // Garantir que o dia tenha dois dígitos
        const formattedYear = `${year[2]}${year[3]}`; // Garantir que o mês tenha dois dígitos
        return `${formattedMonth}/${formattedYear}`;
    };

    const formattedStartDate = formatDate(startDate);
    const formattedEndDate = formatDate(endDate);

    return `${formattedStartDate} - ${formattedEndDate}`;
}
function formatarMoeda(valor) {
    const valorArredondado = Math.floor(valor * 100) / 100;
    return valorArredondado.toLocaleString('pt-BR', {
        style: 'currency',
        currency: 'BRL',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    });
}
export const montarGrafico = (idElemento, label = "Valorização", tipo = 'dias') => {
    const ctx = document.getElementById(idElemento)
    if(!ctx)
        throw `${idElemento} não encontrado!`;

    let labels = [];
    let dados = [];
    if(tipo == "dias") {}
    if(tipo == "meses") {}
    if(tipo == "semestres") {}
    if(tipo == "anos") {}


    new Chart(ctx, {
        type: "line",
        data: {
            labels,
            datasets: [{
                label,
                data: dados,
            }]
        },
        options: {
            scales: {
                y: {beginAtZero: true}
            }
        }
    })
}
const main = async () => {
    const valores_processados = [...document.getElementsByClassName('valor_resultado_info')]
    valores_processados.forEach( div => {
        if(div.id.includes("display", 0)) return
        div.innerText = formatarMoeda(div.innerText)
    })
    const resultado_opcoes = document.getElementById('resultado-opcoes')
    if(!resultado_opcoes) return;
    const containers_resultado = {
        dia:document.getElementById('dias_resultado_processado'),
        mes:document.getElementById('meses_resultado_processado'),
        ano:document.getElementById('anos_resultado_processado'),
        semestre:document.getElementById('semestres_resultado_processado'),
        geral:document.getElementById('resultadoGeralInfoPeriodo'),
    };
    resultado_opcoes.addEventListener('change', (e) => {
        const {value} = e.target
        Object.values(containers_resultado).forEach(container => {
            if(container.classList.contains('hidden')) return;
            container.classList.add("hidden")
        })
        if(!containers_resultado[value].classList.contains('hidden')) return;
        containers_resultado[value].classList.remove("hidden");
    })
}
main();

