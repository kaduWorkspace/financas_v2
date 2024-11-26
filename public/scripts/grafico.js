export const montarGrafico = (idElemento, label = "Valorização", tipo = 'dias') => {
    const ctx = document.getElementById(idElemento)
    if(!ctx)
        throw `${idElemento} não encontrado!`;

    let labels = [];
    let dados = [];
    if(tipo == "dias") {

    }
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
