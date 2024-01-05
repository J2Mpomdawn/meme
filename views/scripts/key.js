var output_count = 0,
    select_count = 0;

function output(e) {
    let div = document.createElement("div");
    div.setAttribute("id", "terminal-output" + (++output_count));
    div.setAttribute("class", "terminal-output-count");

    let current_span = document.createElement("span");
    current_span.setAttribute("class", "terminal-output-current");
    current_span.textContent = document.getElementById("terminal-input-current").textContent;
    div.appendChild(current_span);

    let content_span = document.createElement("span");
    content_span.setAttribute("class", "terminal-output-content");
    content_span.textContent = e.target.value;
    div.appendChild(content_span);

    document.getElementById("terminal-output").appendChild(div);

    e.target.value = "";
    select_count = output_count + 1;
}

function output_cmd(res_s) {
    let div = document.createElement("div");
    div.setAttribute("class", "terminal-output-nocount");

    const res_j = JSON.parse(res_s),
        keys = Object.keys(res_j),
        key_len = keys.map((s) => {return s.length;}).reduce((a, b) => {return Math.max(a, b);});

    for (let i in keys) {
        let content_span = document.createElement("span");
        content_span.setAttribute("class", "terminal-output-content-nocount");
        content_span.innerHTML = keys[i] + "&nbsp;".repeat(key_len - keys[i].length) + "&nbsp;&nbsp;&nbsp;&nbsp;" + res_j[keys[i]];
        div.appendChild(content_span);

        let br = document.createElement("br");
        div.appendChild(br);
    }

    div.removeChild(div.lastElementChild);

    document.getElementById("terminal-output").appendChild(div);
}

function trim_unescape(val) {
    val = val.replace(/^ +| +$/g, "");
    val = val.replace(/^"|"$/g, "");
    val = val.replaceAll('""', '"');

    return val;
}

function get_cmd_parm(val) {
    const vals = val.split("");

    let start = 0,
        str = false,
        dbl = false,
        arr = [];
    
    for (let cur = 0; cur < vals.length; cur++) {
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

        if (!str && vals[cur] == " ") {
            val = trim_unescape(vals.slice(start, cur).join(""));
            arr.push(val);
            start = cur;
        }
    }

    val = trim_unescape(vals.slice(start).join(""));
    arr.push(val);

    arr = arr.filter((v) => v != "");

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

function post(cmd_parm) {
    const ajax = new XMLHttpRequest;
    ajax.open("POST", document.getElementById("terminal-input-current").textContent, true);
    ajax.setRequestHeader("content-type", "application/json");
    ajax.onreadystatechange = () => {
        if (ajax.readyState === 4 && ajax.status === 200) {
            output_cmd(ajax.response);
            ajax.responseText
        }
    };
    ajax.send(cmd_parm);
}

function ws(cmd_parm) {
    //location.origin.replace("https", "wss") + document.getElementById("terminal-input-current").textContent
}

function resize_input() {
    const current_width = document.getElementById("terminal-input-current").offsetWidth + 2;
    document.getElementById("entry").setAttribute("style", "width: calc(100% - " + current_width + "px)");
}

function cmd(e) {
    const cmd_parms = get_cmd_parm(e.target.value);
    const cmd_parm = cmd_parms.parm;

    switch (cmd_parm.app.toLowerCase()) {
        case "cd":
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
                resize_input();
            }
            break;
        case "cd..":
            {
                output(e);

                if (cmd_parm.args[0].length === 0) {
                    let current = document.getElementById("terminal-input-current").textContent;

                    current = current.substring(0, current.lastIndexOf("/"));
                    document.getElementById("terminal-input-current").textContent = current;
                }
                resize_input();
            }
            break;
        case "exit":
            {
                location.href = "/";
            }
            break;
        default:
            {
                switch (cmd_parms.method) {
                    case "--get":
                        {
                            //
                        }
                        break;
                    case "--post":
                        {
                            post(JSON.stringify(cmd_parm));
                            output(e);
                        }
                        break;
                    case "--ws":
                        {
                            ws(JSON.stringify(cmd_parm));
                        }
                        break;
                }
            }
    }
}

document.getElementById("entry").addEventListener("keydown", (e)=>{
    switch (e.key) {
        case "Enter":
            cmd(e);

            break;
        case "ArrowLeft":
            
            break;
        case "ArrowUp":
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