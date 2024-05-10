const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const mongoose = require('mongoose');
const PORT = 3002;

const app = express();

app.use(cors());
app.use(morgan('dev'));
app.use(express.json());

//conectar a la base de datos de mongo
mongoose.connect('mongodb://35.225.155.213:27017/sopes1-p2', { useNewUrlParser: true, useUnifiedTopology: true });
const db = mongoose.connection;

db.on('error', console.error.bind(console, 'Error de conexión a MongoDB:'));
db.once('open', () => {
    console.log('Conexión exitosa a MongoDB');
});

// esquema de la colección logs
const logSchema = new mongoose.Schema({
    name: String,
    album: String,
    year: String,
    rank: String,
    fecha: String,
    hora: String
});

const Log = mongoose.model('Log', logSchema);

app.get('/', (req, res) => {
    res.send('¡Bienvenido a mi API!');
});

app.get('/logs', async (req, res) => {
    try {
        const logs = await Log.find().sort({fecha: -1, hora: -1}).limit(20);
        res.json(logs);
    } catch (error) {
        console.error("Error al obtener los logs:", error);
        res.status(500).json({ error: 'Error al obtener los logs' });
    }
});

app.listen(PORT, () => {
    console.log(`Servidor corriendo en el puerto ${PORT}`);
});