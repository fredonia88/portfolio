fetch('/bls/chart-data')
    .then(response => response.json())
    .then(datasets => {
        const ctx = document.getElementById('blsMedianIncomeChart').getContext('2d');
        new Chart(ctx, {
            type: 'line',
            data: {
                datasets: datasets.map(dataset => ({
                    ...dataset,
                    radius: 6 
                }))
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true, 
                        text: 'Median Income by Demographic Age',
                        color: '#008000',
                        font: {size: 24, family: 'Fira Code'}
                    },
                    legend: {
                        display: true,
                        position: 'top',
                        labels: {
                            font: {size: 14, family: 'Fira Code', style: 'italic'},
                            color: '#008000',
                            boxWidth: 20,
                            boxHeight: 10,
                            padding: 20,
                            usePointStyle: true, 
                            pointStyle: 'circle' 
                        }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                const datasetLabel = context.dataset.label || '';
                                const value = context.raw;
                                const xValue = value.x;
                                const yValue = Math.round(value.y);
                                const formattedYValue = '$' + yValue.toLocaleString()
                                return [
                                    `${datasetLabel}`,
                                    `Year: ${xValue}`, 
                                    `Median Income: ${formattedYValue}`
                                ];
                            }
                        }
                    }
                },
                layout: {
                    padding: 0
                },
                scales: {
                    x: {
                        type: 'category',
                        grid: {
                            color: 'rgba(200, 200, 200, 0.8)', 
                            lineWidth: 0.5,
                        },
                        position: 'bottom', 
                        title: {
                            display: true, 
                            text: 'Year',
                            color: '#008000',
                            font: {size: 20, family: 'Fira Code'}, 
                            padding: {top: 15} 
                        },
                        ticks: {
                            color: '#008000', 
                            font: {size: 16, family: 'Fira Code'},
                            padding: 10
                        }
                    },
                    y: {
                        title: {
                            display: true, 
                            text: 'Median Income',
                            color: '#008000',
                            font: {size: 20, family: 'Fira Code'},
                            padding: {bottom: 15}
                        },
                        grid: {
                            color: 'rgba(200, 200, 200, 0.8)',
                            lineWidth: 0.5,
                        },
                        ticks: {
                            color: '#008000',
                            font: {size: 16, family: 'Fira Code'}, 
                            padding: 10,
                            callback: function(value, index, values) {
                                return '$' + value.toLocaleString();
                            }
                        }
                    }
                }
            }
        });
    })
    .catch(error => console.error('Error fetching chart data:', error));
