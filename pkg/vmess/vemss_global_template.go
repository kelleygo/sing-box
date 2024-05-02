package vmess

var vmessGlobalTemplate = `{
    "log": {
        "disabled": false,
        "level": "error",
        "output": "{{ .logPath }}",
        "timestamp": false
    },
    "dns": {
        "servers": [
            {
                "tag": "dns_proxy",
                "address": "8.8.8.8",
                "address_strategy": "prefer_ipv4",
                "strategy": "prefer_ipv4",
                "detour": "proxy"
            },
            {
                "tag": "dns_local",
                "address": "local",
                "address_strategy": "prefer_ipv4",
                "strategy": "prefer_ipv4",
                "detour": "direct"
            },
            {
                "tag": "dns_block",
                "address": "rcode://refused"
            }
        ],
        "rules": [
            {
                "outbound": "any",
                "server": "dns_local"
            }
        ],
        "final": "dns_proxy",
        "strategy": "prefer_ipv4",
        "disable_cache": false,
        "disable_expire": false,
        "independent_cache": false,
        "reverse_mapping": false,
        "fakeip": {}
    },
    "inbounds": [
        {
            "type": "tun",
            "interface_name": "global_tun",
            "inet4_address": "172.19.0.1/30",
            "mtu": 9000,
            "auto_route": true,
            "strict_route": true,
            "sniff": true
        }
    ],
    "outbounds": [
        {
            "type": "vmess",
            "tag": "proxy",
            "server": "{{ .serviceAddr }}",
            "server_port": {{ .servicePort }},
            "uuid": "{{ .userUuid }}",
            "security": "auto",
            "alter_id": 0,
            "tls": {
            "enabled": true,
            "server_name": "",
            "insecure": true
            }
        },
        {
            "type": "direct",
            "tag": "direct"
        },
        {
            "type": "block",
            "tag": "block"
        },
        {
            "type": "dns",
            "tag": "dns_out"
        }
    ],
    "route": {
        "auto_detect_interface": true,
        "rules": [
            {
                "protocol": "dns",
                "outbound": "dns_out"
            },
            {
                "protocol": ["quic"],
                "outbound": "block"
            },
            {
                "type": "logical",
                "mode": "and",
                "rules": [
                    {
                        "rule_set": "geoip-cn",
                        "invert": true
                    },
                    {
                        "rule_set": "geosite-geolocation-!cn"
                    }
                ],
                "outbound": "proxy"
            },
            {
                "rule_set": "geoip-cn",
                "outbound": "proxy"
            },
            {
                "ip_is_private": true,
                "outbound": "direct"
            }
        ],
        "rule_set": [
            {
                "type": "remote",
                "tag": "geosite-geolocation-!cn",
                "format": "binary",
                "url": "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-geolocation-!cn.srs",
                "download_detour": "proxy"
            },
            {
                "tag": "geosite-cn",
                "type": "remote",
                "format": "binary",
                "url": "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs",
                "download_detour": "proxy"
            },
            {
                "tag": "geoip-cn",
                "type": "remote",
                "format": "binary",
                "url": "https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs",
                "download_detour": "proxy"
            }
        ]
    }
  }`
