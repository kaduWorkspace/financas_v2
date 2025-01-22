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
    let { dados_calculo } = window
    if(!dados_calculo) {
        throw "Nenhum dado para processar!"
    }
    const processar_dados_por_periodo = (dados) => {
        let maior_valorizacao = {valorizacao:0};
        let menor_valorizacao = {valorizacao: -1};
        for (const ano of dados.anos) {
            if(ano.valorizacao > maior_valorizacao.valorizacao)
                maior_valorizacao = ano
            if(menor_valorizacao.valorizacao < 0 || ano.valorizacao < menor_valorizacao.valorizacao)
                menor_valorizacao = ano
        }
        const anos_processado = {
            maior_valorizacao,
            menor_valorizacao
        }
        maior_valorizacao = {valorizacao:0};
        menor_valorizacao = {valorizacao: -1};

        for (const semestre of dados.semestres) {
            if(semestre.valorizacao > maior_valorizacao.valorizacao)
                maior_valorizacao = semestre
            if(menor_valorizacao.valorizacao < 0 || semestre.valorizacao < menor_valorizacao.valorizacao)
                menor_valorizacao = semestre
        }

        const semestres_processado = {
            maior_valorizacao,
            menor_valorizacao
        }
        maior_valorizacao = {valorizacao:0};
        menor_valorizacao = {valorizacao: -1};

        for (const mes of dados.meses) {
            if(mes.valorizacao > maior_valorizacao.valorizacao)
                maior_valorizacao = mes
            if(menor_valorizacao.valorizacao < 0 || mes.valorizacao < menor_valorizacao.valorizacao)
                menor_valorizacao = mes
        }
        const meses_processado = {
            maior_valorizacao,
            menor_valorizacao
        }
        maior_valorizacao = {valorizacao:0};
        menor_valorizacao = {valorizacao: -1};

        for (const dia of dados.dias) {
            if(dia.valorizacao > maior_valorizacao.valorizacao)
                maior_valorizacao = dia
            if(menor_valorizacao.valorizacao < 0 || dia.valorizacao < menor_valorizacao.valorizacao)
                menor_valorizacao = dia
        }
        const dias_processado = {
            maior_valorizacao,
            menor_valorizacao
        }
        return {anos_processado, semestres_processado, meses_processado, dias_processado}
    }
    const template_card_pequeno_resultado = (dados) => {
        return `<div class="flex flex-col w-full items-center bg-gray-700 p-4 rounded-md animate-fade-in">
            <img class="w-8 h-8 mb-2" src="public/images/${dados.icone}" alt="${dados.alt}">
            <span class="text-white text-lg font-semibold">${dados.valor}</span>
            <span class="text-gray-400 text-sm">${dados.titulo}</span>
            </div>`
    }
    const template_resultado = (cards, identificador) => {
        return `
            <div id="${identificador}" class="flex hidden justify-around sm:flex-col flex-row">
            <div class="flex flex-col w-full sm:flex-row sm:mb-2 sm:justify-around gap-4">
            ${cards.shift()}
        ${cards.shift()}
        ${cards.shift()}
            </div>
            <div class="flex flex-col w-full sm:flex-row sm:justify-around gap-4">
            ${cards.shift()}
        ${cards.shift()}
        ${cards.shift()}
            </div>
            </div>`
    }
    const montar_templates = (dados_processados) => {
        const {dias_processado, meses_processado, semestres_processado, anos_processado} = dados_processados
        const cards_dias = [
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( dias_processado.maior_valorizacao.valorizacao ),
                titulo: "Maior valorização",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( dias_processado.maior_valorizacao.resultado_com_valorizacao ),
                titulo: "Valor em caixa",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: dias_processado.maior_valorizacao.data,
                titulo: "Dia",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( dias_processado.menor_valorizacao.valorizacao ),
                titulo: "Menor valorização",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( dias_processado.menor_valorizacao.resultado_com_valorizacao ),
                titulo: "Valor em caixa",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: dias_processado.menor_valorizacao.data,
                titulo: "Dia",
            }),
        ]
        const cards_meses = [
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( meses_processado.maior_valorizacao.valorizacao ),
                titulo: "Maior valorização",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( meses_processado.maior_valorizacao.resultado_com_valorizacao ),
                titulo: "Valor em caixa",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: `${formatarDatasParaTemplate( meses_processado.maior_valorizacao.data_inicial,  meses_processado.maior_valorizacao.data_final )}`,
                titulo: "Data",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( meses_processado.menor_valorizacao.valorizacao ),
                titulo: "Menor valorização",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( meses_processado.menor_valorizacao.resultado_com_valorizacao ),
                titulo: "Valor em caixa",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: `${formatarDatasParaTemplate( meses_processado.menor_valorizacao.data_inicial, meses_processado.menor_valorizacao.data_final )}`,
                titulo: "Data",
            }),
        ]
        const cards_semestres = [
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( semestres_processado.maior_valorizacao.valorizacao ),
                titulo: "Maior valorização",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( semestres_processado.maior_valorizacao.resultado_com_valorizacao ),
                titulo: "Valor em caixa",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: `${formatarDatasParaTemplate( semestres_processado.maior_valorizacao.data_inicial, semestres_processado.maior_valorizacao.data_final )}`,
                titulo: "Data",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( semestres_processado.menor_valorizacao.valorizacao ),
                titulo: "Menor valorização",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( semestres_processado.menor_valorizacao.resultado_com_valorizacao ),
                titulo: "Valor em caixa",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: `${formatarDatasParaTemplate( semestres_processado.menor_valorizacao.data_inicial, semestres_processado.menor_valorizacao.data_final )}`,
                titulo: "Data",
            }),
        ]
        const cards_anos = [
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( anos_processado.maior_valorizacao.valorizacao ),
                titulo: "Maior valorização",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( anos_processado.maior_valorizacao.resultado_com_valorizacao ),
                titulo: "Valor em caixa",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: `${formatarDatasParaTemplate( anos_processado.maior_valorizacao.data_inicial, anos_processado.maior_valorizacao.data_final )}`,
                titulo: "Data",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( anos_processado.menor_valorizacao.valorizacao ),
                titulo: "Menor valorização",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: formatarMoeda( anos_processado.menor_valorizacao.resultado_com_valorizacao ),
                titulo: "Valor em caixa",
            }),
            template_card_pequeno_resultado({
                icone:"icons8-up-64.png",
                alt: "placeholder",
                valor: `${formatarDatasParaTemplate( anos_processado.menor_valorizacao.data_inicial, anos_processado.menor_valorizacao.data_final )}`,
                titulo: "Data",
            }),
        ]
        return [
            template_resultado(cards_dias, "resultado_dias"),
            template_resultado(cards_meses, "resultado_meses"),
            template_resultado(cards_anos, "resultado_anos"),
            template_resultado(cards_semestres, "resultado_semestres"),
        ]
    }
    const resultado_geral_info_container = document.getElementById('resultadoGeralInfo')
    const dados_processados_por_periodo = processar_dados_por_periodo(dados_calculo)
    const templates_resultado = montar_templates(dados_processados_por_periodo)
    resultado_geral_info_container.innerHTML+= templates_resultado.join("")
    console.log("Adiciondo containers de exibição de resultados", resultado_geral_info_container);
    const resultado_opcoes = document.getElementById('resultado-opcoes')
    const containers_resultado = {
        dia:document.getElementById('resultado_dias'),
        mes:document.getElementById('resultado_meses'),
        ano:document.getElementById('resultado_anos'),
        semestre:document.getElementById('resultado_semestres'),
        geral:document.getElementById('resultadoGeralInfoPeriodo'),
    };
    resultado_opcoes.addEventListener('click', (e) => {
        const {value} = e.target
        console.log(value)
        Object.values(containers_resultado).forEach(container => {
            if(container.classList.contains('hidden')) return;
            container.classList.add("hidden")
        })
        if(!containers_resultado[value].classList.contains('hidden')) return;
        containers_resultado[value].classList.remove("hidden");
    })
}
main();

