function loadConfig() {
    query("/openwa/config", "GET", {}, true)
    .done(function (cfg) {
        $("#api_url").val(cfg.api_url || "http://localhost:3000")
        $("#api_key").val(cfg.api_key || "")
        $("#min_delay").val(cfg.min_delay || 3)
        $("#max_delay").val(cfg.max_delay || 8)
    })
}

function saveConfig() {
    var data = {
        api_url: $("#api_url").val(),
        api_key: $("#api_key").val(),
        min_delay: parseInt($("#min_delay").val()) || 3,
        max_delay: parseInt($("#max_delay").val()) || 8
    }
    query("/openwa/config", "POST", data, true)
    .done(function () {
        $("#alert-area").html('<div class="alert alert-success"><i class="fa fa-check"></i> Configurações salvas com sucesso!</div>')
        setTimeout(function () { $("#alert-area").empty() }, 3000)
    })
    .fail(function () {
        $("#alert-area").html('<div class="alert alert-danger">Erro ao salvar configurações.</div>')
    })
}

function testConnection() {
    var url = $("#api_url").val()
    if (!url) {
        $("#alert-area").html('<div class="alert alert-warning">Informe a URL da API OpenWA.</div>')
        return
    }
    $("#alert-area").html('<div class="alert alert-info"><i class="fa fa-spinner fa-spin"></i> Testando conexão...</div>')

    // Save config first then test via a simple GET to the config endpoint as proxy check
    saveConfig()
    query("/openwa/config", "GET", {}, true)
    .done(function () {
        $("#alert-area").html('<div class="alert alert-success"><i class="fa fa-check"></i> Configuração salva. Certifique-se que o OpenWA está rodando em <strong>' + url + '</strong>.</div>')
    })
    .fail(function () {
        $("#alert-area").html('<div class="alert alert-danger">Não foi possível alcançar a API.</div>')
    })
}

$(document).ready(function () {
    loadConfig()
})
