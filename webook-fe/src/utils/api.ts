// 使用库提供的默认配置创建实例

import axios from "axios";

// 此时超时配置的默认值是 `0`
const instance = axios.create({
    baseURL:import.meta.env.BASE_URL
});

// 重写库的超时默认值
// 现在，所有使用此实例的请求都将等待2.5秒，然后才会超时
instance.defaults.timeout =3000;
export default instance