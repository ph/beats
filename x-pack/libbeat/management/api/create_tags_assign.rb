require 'http'
require 'json'
require 'time'

KIBANA_URL = "http://localhost:5601/zhu"
KIBANA_USER = "elastic"
KIBANA_PASSWORD = "changeme"
UUID = JSON.parse(File.read("../../data/meta.json"))["uuid"]


def request()
  tag_name = "hello"
  headers = {
    "kbn-xsrf": "1"
  }


  # url = KIBANA_URL + "/api/beats/tag/" + tag_name

  # payload = {
  #   "color": "#ffffff"
  # }
  # r = HTTP.basic_auth(user: KIBANA_USER, pass: KIBANA_PASSWORD)
  #   .headers(headers)
  #   .put(url, json: payload)

  # if ![200, 201].include?(r.code)
  #   puts "bad response code when creating config block: #{r.code}"
  # end
  # puts r


  blocks = [
    {
      "type" => "output",
      "tag" => tag_name,
      "description" => "amazing 1",
      "last_updated" => Time.now().to_i(),
      "config": {
        "_sub_type" => "elasticsearch",
        "hosts" => ["localhost:9200"],
        "username" => "elastic",
        "password" => KIBANA_PASSWORD,
      }
    },
    # {
    #   "type" => "input",
    #   "description" => "amazing 2",
    #   "tag" => tag_name,
    #   "config": {
    #     "_sub_type" => "log",
    #     "paths" => "/var/log",
    #   }
    # },
    # {
    #   "type" => "input",
    #   "description" => "amazing 3",
    #   "tag" => tag_name,
    #   "config": {
    #     "_sub_type" => "log",
    #     "paths" => "/var/another/log",
    #   }
    # },
  ]

  url = KIBANA_URL + "/api/beats/configurations"

  r = HTTP.basic_auth(user: KIBANA_USER, pass: KIBANA_PASSWORD)
    .headers(headers)
    .put(url, json: blocks)

  if ![200, 201].include?(r.code)
    puts "bad response code when creating config block: #{r.code}"
  end

  puts r
end

request()
