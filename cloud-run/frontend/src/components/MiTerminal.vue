<template>
    <h1>Logs del Sistema de Votaciones</h1>
    <div class="terminal-container">
        <div class="terminal">
            <div class="terminal-consola-side">
                <textarea class="terminal-consola" ref="console" rows="10" v-model="logs"></textarea>
            </div>
            <div class="terminal-boton-side">
                <button class="terminal-boton" @click="actualizarLogs">Actualizar</button>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'MiTerminal',
    data() {
        return {
            logs: ' '
        };
    },
    methods: {
        actualizarLogs() {
            console.log('Actualizando logs...');
            var contador = 0;
            // hacer peticion a http://localhost:3002/logs
            fetch('https://backend-qoj2zdve2q-uk.a.run.app/logs')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Error al obtener los logs');
                    }
                    return response.json();
                })
                .then(data => {
                    const logsFormatted = data.map(log => {
                        contador++;
                        return ` ${contador}. > ${log.fecha} - ${log.hora} - ${log.name} - ${log.album}`;
                    });

                    this.logs = logsFormatted.join('\n');
                })
                .catch(error => {
                    console.error('Error:', error);
                    this.logs = 'Error al obtener los logs';
                });
        }
    }
};
</script>

<style>
    .terminal-container {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
    }

    .terminal{
        display: grid;
        grid-template-columns: 80% 20%;
        width: 100%;
    }

    .terminal-consola-side {
        padding: 10px;
        height: 100vh;
    }

    .terminal-boton-side {
        padding: 10px;
    }

    .terminal-consola {
        width: 100%;
        height: 60%;
        border: 2px solid #000;
        border-radius: 10px;
    }

    .terminal-boton {
        width: 100%;
        height: 10%;
        border: 2px solid #000;
        border-radius: 10px;
        background-color: #000;
        color: #fff;
        font-size: 1.5em;
    }

    .terminal-boton:hover {
        background-color: #fff;
        color: #000;
        cursor: pointer;
    }
</style>