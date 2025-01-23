import {dadosDev} from './dados.js'
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
const STAGE = "DEV"
let wrapper_graficos_dev = ({valorizacao, valor_final, gastos, diferenca, valor_inicial}) => `
    <div id="wrapperGraficos" class="flex flex-col justify-center items-center w-5/6">
                <div class="flex flex-col">
                    <div class="flex flex-col sm:flex-row w-full mb-2 justify-center">
                        <div class="flex flex-col m-1 sm:w-2/4 w-full">
                            <label for="periodo" class="text-gray-400 text-sm mb-1">
                                Escolha o Período
                            </label>
                            <select id="resultado-opcoes" data-ignore_input="1" class="bg-neutral-900 rounded-lg shadow-md p-3 text-gray-300 focus:ring-2 focus:ring-blue-400 focus:outline-none transition-all duration-300">
                                <option value="geral">Resultado Geral</option>
                                <option value="dia">Viualizar por dias</option>
                                <option value="mes">Vizualizar por meses</option>
                                <option value="semestre">Vizualizar por semestres</option>
                                <option value="ano">Vizualizar por anos</option>
                            </select>
                        </div>
                    </div>
                    <div id="resultadoGeralInfo" class="resultado flex-col flex w-full justify-around">
                        <!-- ESTA DIV ABAIXO DEVE SER APAGADA POIS É GERADA NA REQUISICAO -->
                        <div id="resultadoGeralInfoPeriodo" class="flex flex-col w-full justify-center p-4">
                            <div class="grid gap-4 p-4 rounded-lg w-full max-w-4xl grid-cols-1 sm:grid-cols-2 md:grid-cols-4">
                                <div class="flex flex-col items-center bg-gray-700 p-4 rounded-md animate-fade-in">
                                    <img class="w-8 h-8 mb-2" src="public/images/icons8-up-64.png" alt="Valorização">
                                    <span id="display_valorizacao" class="text-white text-lg font-semibold">R$ ${valorizacao}</span>
                                    <span class="text-gray-400 text-sm">Valorização</span>
                                </div>
                                <div class="flex flex-col items-center bg-gray-700 p-4 rounded-md animate-fade-in">
                                    <img class="w-8 h-8 mb-2" src="public/images/icons8-money-with-wings-48.png" alt="Investido">
                                    <span id="display_valor_investido" class="text-white text-lg font-semibold">R$ ${gastos}</span>
                                    <span class="text-gray-400 text-sm">Investido</span>
                                </div>
                                <div class="flex flex-col items-center bg-gray-700 p-4 rounded-md animate-fade-in">
                                    <img class="w-8 h-8 mb-2" src="public/images/icons8-plus-48.png" alt="Lucro">
                                    <span id="diplay_lucro" class="text-white text-lg font-semibold">R$ ${diferenca}</span>
                                    <span class="text-gray-400 text-sm">Lucro</span>
                                </div>
                                <div class="flex flex-col items-center bg-gray-700 p-4 rounded-md animate-fade-in">
                                    <img class="w-8 h-8 mb-2" src="public/images/icons8-race-flag-64.png" alt="Valor inicial">
                                    <span id="display_valor_inicial" class="text-white text-lg font-semibold">R$ ${valor_inicial}</span>
                                    <span class="text-gray-400 text-sm">Valor inicial</span>
                                </div>
                                <div class="flex flex-col items-center bg-green-700 p-6 rounded-md col-span-1 sm:col-span-2 md:col-span-4 animate-fade-in">
                                    <img class="w-10 h-10 mb-3" src="public/images/icons8-money-48.png" alt="Valor final">
                                    <span id="display_valor_final" class="text-white text-2xl font-bold">R$ ${valor_final}</span>
                                    <span class="text-white text-lg">Valor final</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="grafico" class="absolute top-0 bg-preto hidden w-screen overflow-x-auto h-screen">
                    <div class="flex w-full flex-col items-center h-screen">
                        <div class="flex w-full justify-end p-3">
                            <!--<div id="fecharGrafico" class="w-fit cursor-pointer">-->
                            <div id="fecharGrafico" class="absolute top-6 right-6 w-fit cursor-pointer">
                                <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 100 100">
                                    <line x1="10" y1="10" x2="90" y2="90" stroke="white" stroke-width="10" />
                                    <line x1="90" y1="10" x2="10" y2="90" stroke="white" stroke-width="10" />
                                </svg>
                            </div>
                            </div>
                            <div id="resultadoGeral" class="flex flex-col">
                                <h1 class="text-4xl">Resultado geral</h1>
                                <div id="resultadoGeralInfoPeriodo" class="flex w-full justify-around">
                                </div>
                            </div>
                            <h1 class="text-4xl">Tabela</h1>
                            <div class="flex w-full max-height-tabela flex-col items-center" id="resultadoGrafico">

                            </div>
                            <h1 class="text-4xl">Gráfico</h1>
                            <div class="flex w-2/3">
                                <canvas id="chartjs"></canvas>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
 `
if (STAGE == "DEV") {
    const alvo = document.getElementById("alvo_display")
    alvo.innerHTML+= wrapper_graficos_dev(dadosDev)
}

//################ esta parte abaixo deve estar presente no scripts graficos.js ################
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
const dados_processados_por_periodo = processar_dados_por_periodo(dadosDev)
const templates_resultado = montar_templates(dados_processados_por_periodo)
resultado_geral_info_container.innerHTML+= templates_resultado.join("")
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
    console.log(containers_resultado[value].innerHTML)
})
