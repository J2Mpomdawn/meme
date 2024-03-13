var output_count = 0,
    select_count = 0,
    wsc,
    tmp,
    dpr_prev,
    h_f,
    readyState = 4;

/**scroll down*/
function bottom() {
    const e = document.documentElement,
        btm = e.scrollHeight - e.clientHeight;

    window.scroll(0, btm);
}

/**command execution*/
function cmd(e) {
    const cmd_parms = get_cmd_parm(e.target.value);
    const cmd_parm = cmd_parms.parm;

    //distribute commands
    switch (cmd_parms.method) {
        case "--get":
            //get request
            {
                //
            }
            break;
        case "--post":
            //post request
            {
                post(JSON.stringify(cmd_parm));
                output(e);
            }
            break;
        case "--ws":
            //websocket communication
            {
                const app = cmd_parm.app.toLowerCase();
                switch (app) {
                    case "go":
                        //open
                        {
                            ws(cmd_parm.args);
                            output(e);
                        }
                        break;
                    default:
                        //send
                        {
                            wsc.send((new TextEncoder()).encode(app.replaceAll("\\", "\\\\")));
                            output(e);
                        }
                }
            }
            break;
        default:
            //other commands
            {
                switch (cmd_parm.app.toLowerCase()) {
                    case "cd":
                        //move
                        {
                            output(e);

                            const arg = cmd_parm.args[0];
                            let current = document.getElementById("terminal-input-current").textContent;

                            if (arg == "..") {
                                current = current.substring(0, current.lastIndexOf("/"));
                                document.getElementById("terminal-input-current").textContent = current;
                            } else {
                                document.getElementById("terminal-input-current").textContent += "/" + arg;
                            }

                            resize_entry();
                        }
                        break;
                    case "cd..":
                        //one back
                        {
                            output(e);
            
                            let current = document.getElementById("terminal-input-current").textContent;
                            current = current.substring(0, current.lastIndexOf("/"));

                            document.getElementById("terminal-input-current").textContent = current;

                            //entry resizing
                            resize_entry();
                        }
                        break;
                    case "cls":
                        //clear
                        {
                            const toc = document.getElementsByClassName("terminal-output-count");
                            for(let i = 0; i < toc.length; i++) {
                                toc[i].remove();
                            }

                            const ton = document.getElementsByClassName("terminal-output-nocount");
                            for(let i = 0; i < toc.length; i++) {
                                ton[i].remove();
                            }
                        }
                        break;
                    case "exit":
                        //exit key
                        {
                            location.replace("/");
                        }
                        break;
                    case "upload":
                        //register dump file data
                        {
                            upload_dump();
                            output(e);
                        }
                    case "'":
                        //echo
                        {
                            //output(e);
                            output_cmd([cmd_parm.args])
                        }
                        break;
                    default:
                        //invalid command
                        {
                            output(e);

                            output_cmd([
                                "'" + cmd_parm.app +"' is not recognized as an internal or external command,",
                                "operable program or batch file"
                            ]);
                        }
                }
            }
    }

    tmp = cmd_parm;
}

/**analyze pint values*/
function get_cmd_parm(val) {
    const vals = val.split("");

    let start = 0,
        str = false,
        dbl = false,
        arr = [];
    
    //check one letter at a time
    for (let cur = 0; cur < vals.length; cur++) {
        //make the first double quotation mark found the start of the string
        //if the next character after the double quotation mark is not
        //a double quotation mark, make it the end of the string
        if (vals[cur] == '"') {
            if (str && cur != vals.length) {
                if (vals[cur+1] != '"') {
                    if (dbl) {
                        dbl = !dbl;
                    } else {
                        str = !str;
                    }
                } else {
                    dbl = !dbl;
                }
            } else {
                str = !str;
            }
        }

        //if there is a space outside of a string,
        //it is assumed to be a single string up to that point
        if (!str && vals[cur] == " ") {
            //join and unescape
            val = trim_unescape(vals.slice(start, cur).join(""));

            //record string
            arr.push(val);

            //update starting position of a string
            start = cur;
        }
    }

    //join and unescape
    val = trim_unescape(vals.slice(start).join(""));

    //record string
    arr.push(val);

    //excluding empty strings
    arr = arr.filter((v) => v != "");

    //return a struct with the first as the app and the rest as args
    switch (arr[0]) {
        case "--get":
        case "--post":
        case "--ws":
            cmd_parms = {
                "method": arr[0],
                "parm": {
                    "app": arr[1],
                    "args": arr.slice(2)
                }
            };
            break;
        default:
            cmd_parms = {
                "method": "",
                "parm": {
                    "app": arr[0],
                    "args": arr.slice(1)
                }
            };
    }

    return cmd_parms;
}

/**output contents of entry*/
function output(e) {
    //prepare box
    let div = document.createElement("div");
    div.setAttribute("id", "terminal-output" + (++output_count));
    div.setAttribute("class", "terminal-output-count");

    //output content
    let current_span = document.createElement("span");
    current_span.setAttribute("class", "terminal-output-current");
    current_span.textContent = document.getElementById("terminal-input-current").textContent;
    div.appendChild(current_span);

    let content_span = document.createElement("span");
    content_span.setAttribute("class", "terminal-output-content");
    content_span.textContent = e.target.value;
    div.appendChild(content_span);

    document.getElementById("terminal-output").appendChild(div);

    //count and empty entry
    e.target.value = "";
    select_count = output_count + 1;

    //scroll down
    bottom();
}

/**output cmd results*/
function output_cmd(contents) {
    //prepare box
    let div = document.createElement("div");
    div.setAttribute("class", "terminal-output-nocount");

    for (let i = 0; i < contents.length; i++) {
        //output content
        let content_span = document.createElement("span");
        content_span.setAttribute("class", "terminal-output-content-nocount");
        content_span.innerHTML = contents[i];
        div.appendChild(content_span);

        let br = document.createElement("br");
        div.appendChild(br);
    }

    //remove last line break
    div.removeChild(div.lastElementChild);

    //display
    document.getElementById("terminal-output").appendChild(div);

    //scroll down
    bottom();
}

/**output post response*/
function output_cmd_post(res_s) {
    //prepare box
    let div = document.createElement("div");
    div.setAttribute("class", "terminal-output-nocount");

    if (res_s != ""){
        let res_j = JSON.parse(res_s);
        const keys = Object.keys(res_j),
            max_key_len = keys.map((s) => {return s.length;}).reduce((a, b) => {return Math.max(a, b);});

            switch (keys[0]) {
            case "failed":
                //failure of some kind
                {
                    output_cmd(res_j[keys[0]].join("\n"));
                }
                break;
            case "ws":
                //make websocket
                {
                    ws(tmp.args);
                }
                break;
            case "table":
                //output select results
                {
                    let cols,
                        col_len = new Map(),
                        row_line;

                    res_j = JSON.parse(res_j.table);

                    //measure the maximum width of each column
                    //check one line at a time
                    for (let i = 0; i < res_j.length; i++) {
                        //processing for the first line only
                        if (i == 0) {
                            //get all column names
                            cols = Object.keys(res_j[0])

                            //column direction
                            for (let j = 0; j < cols.length; j++) {
                                //column name width
                                col_len.set(cols[j], str_width(cols[j]).width);

                                //row direction
                                for (let k = 0; k < res_j.length; k++) {
                                    //value width
                                    const val_len = str_width(res_j[k][cols[j]]).width;

                                    //if this one is longer, it will be overwritten
                                    if(val_len > col_len.get(cols[j])) {
                                        col_len.set(cols[j], val_len);
                                    }
                                }
                            }
                        }

                        let row;
                        //processing for the first line only
                        if (i == 0) {
                            //prepare a one-line border
                            row = "+";
                            for (let j = 0; j < cols.length; j++) {
                                row += "-".repeat(col_len.get(cols[j])) + "+";
                            }
                            row_line = row;

                            //draw a border
                            let content_span = document.createElement("span");
                            content_span.setAttribute("class", "terminal-output-content-nocount");
                            content_span.innerHTML = row_line;
                            div.appendChild(content_span);
            
                            let br = document.createElement("br");
                            div.appendChild(br);

                            //draw column names, taking into account the maximum width per column
                            row = "|";
                            for (let j = 0; j < cols.length; j++) {
                                const val = cols[j];
                                row += val + "　".repeat(str_width(val).full) + "&nbsp;".repeat(col_len.get(cols[j]) - str_width(val).width) + "|";
                            }
                            content_span = document.createElement("span");
                            content_span.setAttribute("class", "terminal-output-content-nocount");
                            content_span.innerHTML = row;
                            div.appendChild(content_span);

                            br = document.createElement("br");
                            div.appendChild(br);

                            //draw a border
                            content_span = document.createElement("span");
                            content_span.setAttribute("class", "terminal-output-content-nocount");
                            content_span.innerHTML = row_line;
                            div.appendChild(content_span);

                            br = document.createElement("br");
                            div.appendChild(br);
                        }

                        //draw value, taking into account the maximum width per column
                        row = "|";
                        for (let j = 0; j < cols.length; j++) {
                            const val = res_j[i][cols[j]];
                            row += val + "　".repeat(str_width(val).full) + "&nbsp;".repeat(col_len.get(cols[j]) - str_width(val).width) + "|";
                        }
                        content_span = document.createElement("span");
                        content_span.setAttribute("class", "terminal-output-content-nocount");
                        content_span.innerHTML = row;
                        div.appendChild(content_span);

                        let br = document.createElement("br");
                        div.appendChild(br);

                        //draw a border
                        content_span = document.createElement("span");
                        content_span.setAttribute("class", "terminal-output-content-nocount");
                        content_span.innerHTML = row_line;
                        div.appendChild(content_span);

                        br = document.createElement("br");
                        div.appendChild(br);
                    }

                    //remove last line break
                    div.removeChild(div.lastElementChild);

                    div.classList.add("terminal-sql-table");

                    //display
                    document.getElementById("terminal-output").appendChild(div);
                }
                break;
            default:
                {
                    //display response with key and value set
                    for (let i = 0; i < keys.length; i++) {
                        let content_span = document.createElement("span");
                        content_span.setAttribute("class", "terminal-output-content-nocount");
                        content_span.innerHTML = keys[i] + "&nbsp;".repeat(max_key_len - keys[i].length) + "&nbsp;&nbsp;&nbsp;&nbsp;" + res_j[keys[i]];
                        div.appendChild(content_span);
        
                        let br = document.createElement("br");
                        div.appendChild(br);
                    }

                    //remove last line break
                    div.removeChild(div.lastElementChild);

                    //display
                    document.getElementById("terminal-output").appendChild(div);
                }
        }

        //scroll down
        bottom();
    }
}

/**post request*/
function post(cmd_parm) {
    //prepare for communicaion
    const ajax = new XMLHttpRequest;
    ajax.open("POST", document.getElementById("terminal-input-current").textContent, true);
    ajax.setRequestHeader("content-type", "application/json");
    ajax.onreadystatechange = () => {
        if (ajax.readyState === 2) {
            readyState = 2;
        } else if (ajax.readyState === 4) {
            readyState = 4;

            switch (ajax.status) {
                case 200:
                    //ok
                    {
                        output_cmd_post(ajax.response);
                    }
                    break;
                case 404:
                    //mirage
                    {
                        output_cmd_post(ajax.response);
                    }
                    break;
                case 500:
                    //server error
                    {
                    }
                    break;
            }

            //return hidden items when communication ends
            document.getElementById("terminal-input").style.display = "";
            document.getElementById("terminal-input").removeAttribute("disabled");
            document.getElementById("entry").focus();
        }
    };

    //request
    ajax.send(cmd_parm);

    //hide entry during communication
    document.getElementById("terminal-input").style.display = "none";
    document.getElementById("terminal-input").setAttribute("disabled", true);
}

/**entry resizing*/
function resize_entry() {
    //allcate the remaining width to entry
    const current_width = Math.ceil(document.getElementById("terminal-input-current").getBoundingClientRect().width) + 1;
    document.getElementById("entry").setAttribute("style", "width: calc(100% - " + current_width + "px)");
}

/**ratio of half and full width*/
function set_h_f() {
    //measure width of one half character and one full character
    let h = document.getElementById("half").getBoundingClientRect().width,
        f = document.getElementById("full").getBoundingClientRect().width,
        ratio = [];

    //ratio of half and full for 5 patterns
    ratio[0] = f / h;
    ratio[1] = f * 2 / h;
    ratio[2] = f * 3 / h;
    ratio[3] = f * 4 / h;
    ratio[4] = f * 5 / h;

    //decimal portion
    const ratio_dec = ratio.map((x) => {return x - Math.floor(x)});

    //index with the smallest decimal part
    f = ratio_dec.indexOf(Math.min(...ratio_dec));

    //truncate corresponding width to an integer
    h = Math.floor(ratio[f]);

    h_f = {
        "h": h,
        "f": f + 1
    };
}

/**character width and number of full characters*/
function str_width(str) {
    const half = (str.match(/[\x01-\x7E\xA1-\xDF]/g) || []).length,
        full = (str.match(/[^\x01-\x7E\xA1-\xDF]/g) || []).length;

    //calculate width using ratios
    return {
        "full": Math.ceil(full / h_f.f) * h_f.f - full,
        "width": half + Math.ceil(full / h_f.f) * h_f.h
    };
}

/**trim and unescape*/
function trim_unescape(val) {
    //first and last half space
    val = val.replace(/^ +| +$/g, "");

    //first and last double quotation
    val = val.replace(/^"|"$/g, "");

    //two double quotations
    val = val.replaceAll('""', '"');

    return val;
}

/**upload record dump file */
function upload_dump() {
    const input = document.createElement("input");
    input.setAttribute("type", "file");
    input.setAttribute("id", "file");

    //event when file is selected
    input.addEventListener("change", (e) => {
        const file   = e.target.files,
              reader = new FileReader();

        //get file contents
        reader.readAsText(file[0]);

        //event when file has been read
        reader.onload = () => {
            let vals  = "",
                count = 0,
                n     = "1";

            const res  = reader.result,
                  //wait for ms
                  wait = async (ms) => new Promise(resolve => setTimeout(resolve, ms)),
                  //uploading process
                  upld = async () => {
                    //execute if the previous post process is finished
                    if (readyState === 4) {
                        post(JSON.stringify({
                            "app": "upload",
                            "args": [
                                vals.slice(0, -1),
                                n
                            ]
                        }));
                        
                        return;
                    //if not, wait
                    } else {
                        await wait(1000);
                        upld();
                    }
                }

            //get lines
            const rows = res.split("\r\n");
            for(let i = 1; i < rows.length; i++) {
                //execute every 1000 lines
                if(Math.floor(i/1000) !== count) {
                    n += "-" + (i - 1);
                    upld();

                    count = Math.floor(i/1000);
                    vals = "";
                    n    = i + "";
                }

                if(rows[i] == "") {
                    continue;
                }

                let row = "(";
                n++;

                //get values
                const cols = rows[i].split("\t");
                for(let j = 0; j < cols.length; j++) {

                    let cell = cols[j];

                    //delete the first and last '"'
                    if(cell.slice(0, 1) == '"') {
                        cell = cell.slice(1);
                    }
                    if(cell.slice(-1) == '"') {
                        cell = cell.slice(0, -1);
                    }

                    //unescape escaped '"'
                    cell = cell.replaceAll('""', '"');

                    row += "'" + cell + "',";

                }

                row = row.slice(0, -1) + ")";
                vals += row + ",";

            }

            //process remaining data
            n += "-";
            upld();
        }

        //to check file name
        console.log(file[0].name);
    });

    //place object for file processing,
    //delete after click event execution
    const ti = document.getElementById("terminal-input");
    ti.appendChild(input);

    input.click();

    ti.removeChild(ti.lastElementChild);
}

/**websocket*/
function ws(args) {
    //specify path
    let path = ""
    const i = args.indexOf("-u");
    if (i !== -1 && i + 1 < args.length) {
        path = args[i + 1];
    } else {
        path = document.getElementById("terminal-input-current").textContent;
    }

    //connect websocket
    wsc = new WebSocket(location.origin.replace("https", "wss") + path)
    wsc.binaryType = "arraybuffer";

    //when open
    wsc.addEventListener("open", (e) => {
        //indicates communication in progress
        document.getElementById("terminal-input").style.color = "blue";
    });

    //when receive
    wsc.addEventListener("message", (e) => {
        const e_s = String.fromCharCode.apply("", new Uint8Array(e.data));

        //how response
        output_cmd([
            e_s
        ]);
        
        //terminate communication
        if (e_s == "closed") {
            wsc.close();
        }
    });

    //when close
    wsc.addEventListener("close", (e) => {
        //exit display during communicaiton
        document.getElementById("terminal-input").style.color = "";
    });
}

/**click event*/
document.addEventListener("click", () => {
    //focus on entry
    document.getElementById("entry").focus();
});

/**page load completion event*/
document.addEventListener("DOMContentLoaded", () => {
    //magnification power record
    dpr_prev = window.devicePixelRatio;

    //ratio of half and full width
    set_h_f();
});

/**keydown event for entry*/
document.getElementById("entry").addEventListener("keydown", (e) => {
    switch (e.key) {
        case "Enter":
            //command execution
            cmd(e);

            //scroll down
            bottom();

            break;
        case "ArrowLeft":
            
            break;
        case "ArrowUp":
            //input value of previous line
            if (select_count > 1) {
                let log = document.getElementById("terminal-output" + (--select_count)).lastElementChild.textContent;

                e.target.value = log;

                e.target.focus();
                e.target.setSelectionRange(0, log.length);
            }

            break;
        case "ArrowRight":
            
            break;
        case "ArrowDown":
            //input value of next line
            if (select_count < output_count) {
                let log = document.getElementById("terminal-output" + (++select_count)).lastElementChild.textContent;

                e.target.value = log;

                e.target.focus();
                e.target.setSelectionRange(0, log.length);
            } else if (select_count === output_count){
                e.target.value = "";

                select_count++;
            }

            break;
    }
});

/**resize event*/
window.addEventListener("resize", () => {
    //magnification rate has changed
    if (dpr_prev != window.devicePixelRatio) {
        //ratio of half and full width
        set_h_f();
    }

    dpr_prev = window.devicePixelRatio;
});