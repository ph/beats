require 'http'
require 'json'

KIBANA_URL = "http://localhost:5601/zhu"
KIBANA_USER = "elastic"
KIBANA_PASSWORD = "changeme"
UUID = JSON.parse(File.read("../../data/meta.json"))["uuid"]


def request()
  tag_name = "test"
  headers = {
    "kbn-xsrf": "1"
  }

  blocks = {
    "type": "output",
    "configs": [
      {
        "output": "elasticsearch",
        "elasticsearch": {
          "hosts": ["localhost:9200"],
          "username":"elastic",
          "password": KIBANA_PASSWORD,
        }
      }
    ]
  }

  url = KIBANA_URL + "/api/beats/tag" + tag_name
  data = {
    "color": "#DD0A73",
    "configuration_blocks": blocks,
  }

  r = HTTP.basic_auth(user: KIBANA_USER, password: KIBANA_PASSWORD)
    .headers(headers)
    .post(url, json: data)

  if ![200, 201].include?(r.code)
    puts "bad response code when creating config block: ${r.code}"
    puts r
  end


  data = {"assignments": [{"beatId": UUID, "tag": tag_name}]}
  url = KIBANA_URL + "/api/beats/agents_tags/assignments"

  r = HTTP.basic_auth(user: KIBANA_USER, password: KIBANA_PASSWORD)
    .headers(headers)
    .post(url, json: data)

  if ![200, 201].include?(r.code)
    puts "bad response code when assignment tag to beat: ${r.code}"
    puts r
  end
end
