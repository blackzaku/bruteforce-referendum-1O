const fs = require('fs');
const request = require('request');
const crypto = require('crypto');
const DIRSHA = 2;
const FILESHA = 2;
const BUCLE = 1714;
let tried = {};

const BARCELONA = require('./zip_codes.js').BARCELONA;
const CC_A = 65;
const CC_Z = 90;
const MIN_YEAR = 1965;
const MAX_YEAR = 2000;

// let dniN = -1;
// let dniCc = CC_Z;
let dniN = 92612;
let dniCc = 'S';

let zipI = -1;

function resetZip() {
    zipI = -1;
}

let day = 31;
let month = 12;
let year = MIN_YEAR;
function resetDate() {
    day = 31;
    month = 12;
    year = MIN_YEAR;
}

function nextDni() {
    if (dniCc < CC_Z) {
        dniCc += 1;
    } else if (dniN < 99999) {
        dniCc = CC_A;
        dniN += 1;
    } else {
        return null;
    }
    return `0000${dniN}`.slice(-5) + String.fromCharCode(dniCc);
}

function nextZip() {
    if (zipI < BARCELONA.length) {
        zipI += 1;
        return BARCELONA[zipI];
    } else {
        return null;
    }
}

function nextDate() {
    day += 1;
    if (day === 32) {
        month += 1;
        day = 1;
    }
    if (month === 13) {
        year += 1;
        month = 1;
    }
    if (year === MAX_YEAR) {
        return null;
    } else {
        return `${year}${month >= 10 ? month : '0' + month}${day >= 10 ? day : '0' + day}`;
    }
}


function decrypt(text, password){
    const decipher = crypto.createDecipher('aes-256-cbc',password);
    let dec = decipher.update(text,'hex','utf8');
    dec += decipher.final('utf8');
    return dec;
}

function hash(text) {
    return crypto.createHash('sha256').update(text).digest('hex');
}
function bucleHash(key, n) {
    let clauTemp = key;
    for(let x = 0; x < n; x += 1){
        clauTemp = hash(clauTemp);
    }
    return clauTemp;
}

async function check(dni, date, zip) {
    const key = dni + date + zip;
    const firstSha256 = hash(bucleHash(key,BUCLE));
    const secondSha256 = hash(firstSha256);
    const dir = secondSha256.substring(0,DIRSHA);
    const file = secondSha256.substring(DIRSHA,DIRSHA+FILESHA);
    const localFile = `db/${dir}_${file}.db`;
    let lines = tried[localFile];
    if (tried[localFile] === undefined) {
        if (fs.existsSync(localFile)) {
            lines = fs.readFileSync(localFile).toString().split('\n');
            tried[localFile] = lines;
        } else {
            const url = `https://wikileaks.org/mirrors/catref/db/${dir}/${file}.db`;
            lines = await new Promise((resolve, reject) => {
                console.warn(`Requesting file ${localFile}`);
                request(url, function (error, response, body) {
                    if (error || body.substring(0, 12) === 'ipfs resolve' || body.length < 256) {
                        console.error(`ERROR | File ${localFile} FAILED`);
                        tried[localFile] = null;
                        resolve(null);
                    } else {
                        console.warn(`Got file ${localFile}`);
                        fs.writeFileSync(localFile, body);
                        const split = body.toString().split('\n');
                        tried[localFile] = split;
                        resolve(split);
                    }
                });
            });
        }
    }
    if (lines === null) {
        return null
    } else {
        for (let line of lines) {
            if (line.substring(0, 60) === secondSha256.substring(4)) {
                return decrypt(line.substring(60), firstSha256).split('#');
            }
        }
        return null;
    }
}

async function runDniDate(dni, date) {
    resetZip();
    let zip = nextZip();
    while (zip !== null) {
        const info = await check(dni, date, zip);
        if (info !== null) {
            console.log(`
                    #################################
                    # FOUND !
                    # DNI: ${dni}
                    # DATE: ${date}
                    # ZIP: ${zip}
                    # INFO: ${info}
                    #################################`)
        }
        zip = nextZip();
    }
}

async function runDni(dni) {
    resetDate();
    let date = nextDate();
    while (date !== null) {
        await runDniDate(dni, date);
        date = nextDate();
    }
}

async function run() {
    let dni = nextDni();
    while (dni !== null) {
        console.log(`Checking ${dni}`);
        await runDni(dni);
        dni = nextDni();
    }
}

if (fs.existsSync('db')) {
    // Do something
    fs.mkdirSync('db');
}
run().then( () => console.log('DONE!'), error => console.error(error));
