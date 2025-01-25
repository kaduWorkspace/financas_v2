function formatarMoeda(valor) {
    const valorArredondado = Math.floor(valor * 100) / 100;
    return valorArredondado.toLocaleString('pt-BR', {
        style: 'currency',
        currency: 'BRL',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    });
}
export const montarGrafico = (idElemento, label = "Valorização", dados) => {
    const ctx = document.getElementById(idElemento)
    if(!ctx)
        throw `${idElemento} não encontrado!`;
    else if(window.grafico_canva)
        window.grafico_canva.destroy()

    const labels = dados.y
    const data = dados.x
    window.grafico_canva = new Chart(ctx, {
        type: "line",
        data: {
            labels,
            datasets: [{
                label,
                data: data,
            }]
        },
        options: {
            scales: {
                y: {beginAtZero: true}
            }
        }
    })
}
const ajustar_dados_grafico = () => {
    const MAX_ITENS_GRAFICO = 15
    const dados_grafico = {
        dias: {
            dados: {x:[],y:[]},
        },
        meses: {
            dados: {x:[],y:[]},
        },
        semestres: {
            dados: {x:[],y:[]},
        },
        anos: {
            dados: {x:[],y:[]},
        },
        geral: {
            dados: {x:[],y:[]},
        },
    }
    const {anos, dias, meses, semestres} = window.dados_calculo
    const {dados_calculo} = window
    const grafico_dados_geral = meses
    const jump = grafico_dados_geral.length < MAX_ITENS_GRAFICO ? 1 : Math.floor(grafico_dados_geral.length/MAX_ITENS_GRAFICO)
    let index = 0;
    while(index < grafico_dados_geral.length - 1) {
        if(index > grafico_dados_geral.length - 1) {
            const item = grafico_dados_geral.pop();
            dados_grafico.geral.dados.x.push(item.resultado_com_valorizacao)
            dados_grafico.geral.dados.y.push(item.data ? item.data : item.data_inicial)
            break;
        } else {
            dados_grafico.geral.dados.x.push(grafico_dados_geral[index].resultado_com_valorizacao)
            dados_grafico.geral.dados.y.push(grafico_dados_geral[index].data ? grafico_dados_geral[index].data : grafico_dados_geral[index].data_inicial)
        }
        index = index + jump
    }
    Object.keys({dias, meses, semestres, anos}).forEach(key => {
        let index = 0;
        const dados_grafico_periodo = dados_calculo[key]
        const jump = dados_grafico_periodo.length > MAX_ITENS_GRAFICO ? Math.floor(dados_grafico_periodo.length/MAX_ITENS_GRAFICO) : 1 ;
        while(index < dados_grafico_periodo.length - 1) {
            if(index > dados_grafico_periodo.length - 1) {
                const item = dados_grafico_periodo.pop();
                dados_grafico[key].dados.x.push(item.valorizacao)
                dados_grafico[key].dados.y.push(item.data ? item.data : item.data_inicial)
                return;
            } else {
                dados_grafico[key].dados.x.push(dados_grafico_periodo[index].valorizacao)
                dados_grafico[key].dados.y.push(dados_grafico_periodo[index].data ? dados_grafico_periodo[index].data : dados_grafico_periodo[index].data_inicial)
            }
            index = index + jump
        }
    })
    return dados_grafico
}
const main = async () => {
    let dados_grafico;
    if(window?.dados_calculo) {
        dados_grafico = ajustar_dados_grafico();
        const grafico_botao_abrir  = document.getElementById("botao_ativar_grafico");
        const grafico_container = document.getElementById("grafico_container");
        const grafico_botao_fechar = document.getElementById("grafico_fechar");
        grafico_botao_fechar.addEventListener("click", () => {
            grafico_container.classList.add("hidden")
            habilitarScroll();
        })
        grafico_botao_abrir.addEventListener("click", () => {
            grafico_container.classList.remove("hidden")
            moverParaTopo();
            desabilitarScroll();
        })
        montarGrafico('chartjs', 'Valorizacao', dados_grafico.geral.dados)
    }
    const valores_processados = [...document.getElementsByClassName('valor_resultado_info')]
    valores_processados.forEach( div => {
        if(div.id.includes("display", 0)) return
        div.innerText = formatarMoeda(div.innerText)
    })
    const resultado_opcoes = document.getElementById('resultado-opcoes')
    if(!resultado_opcoes) return;
    const containers_resultado = {
        dias:document.getElementById('dias_resultado_processado'),
        meses:document.getElementById('meses_resultado_processado'),
        anos:document.getElementById('anos_resultado_processado'),
        semestres:document.getElementById('semestres_resultado_processado'),
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
        montarGrafico('chartjs', 'Valorizacao', dados_grafico[value].dados)
    })
}
main();

function desabilitarScroll() {
    document.body.style.overflow = 'hidden';
}
function habilitarScroll() {
    document.body.style.overflow = 'auto';
}
function moverParaTopo() {
    window.scrollTo({
        top: 0,
        behavior: 'smooth', // Adiciona animação suave
    });
}
