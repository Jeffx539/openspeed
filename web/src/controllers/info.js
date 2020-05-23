import axios from "axios";

class InfoController {
  static async queryInfo(endpoint) {
    let data = await axios.get(`${endpoint}/info`);
    return data.data;
  }
}

export default InfoController;
