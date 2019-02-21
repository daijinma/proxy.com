module.exports = {
    apps: [
        {
            name: "proxy.com",
            instances: 1,
            script:"./proxy.com",
            exec_mode: "fork",    // 一定要是fork
            interpreter: "./proxy.com",   // windows下加.exe
            env: {             // 环境变量
                env: "development",
            },
        }
    ]
}