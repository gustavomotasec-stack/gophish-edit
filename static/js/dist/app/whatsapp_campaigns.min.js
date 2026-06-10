var waTemplates = []
var phoneLists = []
var goCampaigns = []

function waFormat(text) {
    return text
        .replace(/\*(.*?)\*/g, '<strong>$1</strong>')
        .replace(/_(.*?)_/g, '<em>$1</em>')
        .replace(/~(.*?)~/g, '<del>$1</del>')
        .replace(/```(.*?)```/gs, '<code>$1</code>')
        .replace(/\n/g, '<br>')
}

function loadAll() {
    query("/whatsapp_templates/", "GET", {}, true).done(function (data) {
        waTemplates = data || []
        $.each(waTemplates, function (i, t) {
            $("#template_id").append('<option value="' + t.id + '">' + escapeHtml(t.name) + '</option>')
        })
    })
    query("/phone_lists/", "GET", {}, true).done(function (data) {
        phoneLists = data || []
        $.each(phoneLists, function (i, p) {
            var count = p.numbers ? p.numbers.length : 0
            $("#phone_list_id").append('<option value="' + p.id + '">' + escapeHtml(p.name) + ' (' + count + ' números)</option>')
        })
    })
    query("/campaigns/", "GET", {}, true).done(function (data) {
        goCampaigns = data || []
        $.each(goCampaigns, function (i, c) {
            $("#campaign_id").append('<option value="' + c.id + '">' + escapeHtml(c.name) + '</option>')
        })
    })
}

$("#template_id").on("change", function () {
    var id = parseInt($(this).val())
    var t = waTemplates.find(function (x) { return x.id === id })
    if (!t) { $("#preview-area").hide(); return }
    var preview = t.body
        .replace(/{{\.URL}}/g, 'https://exemplo.com/r/AbCdEfG')
        .replace(/{{\.FirstName}}/g, 'João')
        .replace(/{{\.LastName}}/g, 'Silva')
    $("#preview-box").html(waFormat(preview))
    $("#preview-area").show()
})

$("#phone_list_id").on("change", function () {
    var id = parseInt($(this).val())
    var p = phoneLists.find(function (x) { return x.id === id })
    if (!p) { $("#list-info").text(""); return }
    var count = p.numbers ? p.numbers.length : 0
    $("#list-info").text(count + " números nesta lista")
})

function dispatchCampaign() {
    var templateId = parseInt($("#template_id").val())
    var phoneListId = parseInt($("#phone_list_id").val())
    var campaignId = parseInt($("#campaign_id").val()) || 0
    var minDelay = parseInt($("#min_delay").val()) || 3
    var maxDelay = parseInt($("#max_delay").val()) || 8
    var scheduledAt = $("#scheduled_at").val()

    if (!templateId) {
        $("#modal-alert").html('<div class="alert alert-danger">Selecione um template.</div>')
        return
    }
    if (!phoneListId) {
        $("#modal-alert").html('<div class="alert alert-danger">Selecione uma lista de números.</div>')
        return
    }

    var payload = {
        template_id: templateId,
        phone_list_id: phoneListId,
        campaign_id: campaignId,
        min_delay: minDelay,
        max_delay: maxDelay
    }
    if (scheduledAt) {
        payload.scheduled_at = new Date(scheduledAt).toISOString()
    }

    // Save delay config before dispatch
    query("/openwa/config", "POST", { min_delay: minDelay, max_delay: maxDelay }, true)

    query("/openwa/dispatch", "POST", payload, true)
    .done(function (data) {
        $("#modal").modal("hide")
        Swal.fire({
            title: "Disparo iniciado!",
            text: data.message,
            type: "success",
            timer: 3000,
            showConfirmButton: false
        })
    })
    .fail(function (xhr) {
        var msg = xhr.responseJSON ? xhr.responseJSON.message : "Erro ao disparar."
        $("#modal-alert").html('<div class="alert alert-danger">' + escapeHtml(msg) + '</div>')
    })
}

$("#modal").on("show.bs.modal", function () {
    $("#modal-alert").empty()
    $("#template_id").val("")
    $("#phone_list_id").val("")
    $("#campaign_id").val("0")
    $("#preview-area").hide()
    $("#list-info").text("")
    $("#scheduled_at").val("")
})

$(document).ready(function () {
    loadAll()
    $("#emptyMessage").show()
})

function escapeHtml(s) {
    return String(s)
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
}
