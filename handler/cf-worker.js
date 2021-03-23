// 以 网址?code=获取到的code 来访问

// github应用信息
const client_name = "应用名"
const client_id = "应用id"
const client_secret = "应用秘钥"


addEventListener("fetch", event => {
    return event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
    return start(request)
}

// 开始
async function start(request) {
    const data = await getCode(request)
    return new Response(JSON.stringify(data), {
        headers: {
            "content-type": "application/json;charset=UTF-8",
        },
    })
}

// 从网址中提取code
async function getCode(request) {
    const s = new URL(request.url).searchParams
    if (s.has("code")) {
        code = s.get("code")
        if (code) {
            return await getAccessToken(code)
        }
    }
    return { error: "no code" }
}

// 请求githubAPI，获取access_token
async function getAccessToken(code) {
    const option = {
        method: "POST",
        headers: {
            "accept": "application/json",
        }
    }
    const res = await fetch(`https://github.com/login/oauth/access_token?client_id=${client_id}&client_secret=${client_secret}&code=${code}`, option)
    const j = await res.json()

    if (j.error) {
        return { error: `[${j.error}] ${j.error_description}` }
    }
    return await getUserInfo(j.access_token, j.token_type)
}

// 获取用户数据
async function getUserInfo(access_token, token_type) {
    const option = {
        headers: {
            "Authorization": `${token_type} ${access_token}`,
            "User-Agent": client_name,
            "Accept": "application/json"
        }
    }
    const res = await fetch("https://api.github.com/user", option)
    const data = await res.json()
    return data
}


