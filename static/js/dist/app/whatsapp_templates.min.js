var templates = []
var deleteId = 0

// Formata o texto do WhatsApp para preview no browser
function waFormat(text) {
    return text
        .replace(/\*(.*?)\*/g, '<strong>$1</strong>')
        .replace(/_(.*?)_/g, '<em>$1</em>')
        .replace(/~(.*?)~/g, '<del>$1</del>')
        .replace(/```(.*?)```/gs, '<code>$1</code>')
        .replace(/\n/g, '<br>')
}

function updatePreview() {
    var text = $("#body").val()
        .replace(/{{\.URL}}/g, 'https://exemplo.com/r/AbCdEfG')
        .replace(/{{\.FirstName}}/g, 'João')
        .replace(/{{\.LastName}}/g, 'Silva')
    $("#preview-box").html(waFormat(text))
}

function wrapText(before, after) {
    var ta = document.getElementById("body")
    var start = ta.selectionStart
    var end = ta.selectionEnd
    var val = ta.value
    ta.value = val.substring(0, start) + before + val.substring(start, end) + after + val.substring(end)
    ta.selectionStart = start + before.length
    ta.selectionEnd = end + before.length
    ta.focus()
    updatePreview()
}

function insertText(text) {
    var ta = document.getElementById("body")
    var start = ta.selectionStart
    var val = ta.value
    ta.value = val.substring(0, start) + text + val.substring(start)
    ta.selectionStart = ta.selectionEnd = start + text.length
    ta.focus()
    updatePreview()
}

function loadTemplates() {
    query("/whatsapp_templates/", "GET", {}, true)
    .done(function (data) {
        templates = data || []
        renderTable()
    })
}

function renderTable() {
    var tbody = $("#templateTable tbody").empty()
    if (!templates.length) {
        $("#emptyMessage").show()
        $("#templateTable").hide()
        return
    }
    $("#emptyMessage").hide()
    $("#templateTable").show()
    $.each(templates, function (i, t) {
        var date = moment(t.modified_date).format("DD/MM/YYYY HH:mm")
        var row = '<tr>' +
            '<td>' + escapeHtml(t.name) + '</td>' +
            '<td>' + date + '</td>' +
            '<td>' +
            '<button class="btn btn-xs btn-default" onclick="editTemplate(' + t.id + ')"><i class="fa fa-pencil"></i> Editar</button> ' +
            '<button class="btn btn-xs btn-danger" onclick="confirmDelete(' + t.id + ')"><i class="fa fa-trash"></i> Excluir</button>' +
            '</td></tr>'
        tbody.append(row)
    })
}

function resetModal() {
    $("#template-id").val("0")
    $("#modal-title").text("Novo WhatsApp Template")
    $("#name").val("")
    $("#body").val("")
    $("#preview-box").html("")
    $("#modal-alert").empty()
}

function editTemplate(id) {
    var t = templates.find(function (x) { return x.id === id })
    if (!t) return
    $("#template-id").val(t.id)
    $("#modal-title").text("Editar Template")
    $("#name").val(t.name)
    $("#body").val(t.body)
    updatePreview()
    $("#modal").modal("show")
}

function saveTemplate() {
    var id = parseInt($("#template-id").val())
    var data = { name: $("#name").val(), body: $("#body").val() }
    var method = id === 0 ? "POST" : "PUT"
    var endpoint = id === 0 ? "/whatsapp_templates/" : "/whatsapp_templates/" + id

    query(endpoint, method, data, true)
    .done(function (t) {
        if (id === 0) {
            templates.push(t)
        } else {
            templates = templates.map(function (x) { return x.id === id ? t : x })
        }
        renderTable()
        $("#modal").modal("hide")
    })
    .fail(function (xhr) {
        var msg = xhr.responseJSON ? xhr.responseJSON.message : "Erro ao salvar."
        $("#modal-alert").html('<div class="alert alert-danger">' + escapeHtml(msg) + '</div>')
    })
}

function confirmDelete(id) {
    deleteId = id
    $("#deleteModal").modal("show")
}

$("#confirm-delete").on("click", function () {
    query("/whatsapp_templates/" + deleteId, "DELETE", {}, true)
    .done(function () {
        templates = templates.filter(function (x) { return x.id !== deleteId })
        renderTable()
        $("#deleteModal").modal("hide")
    })
})

// Preview em tempo real
$("#body").on("input", updatePreview)

// Reset ao abrir modal novo
$("#modal").on("show.bs.modal", function (e) {
    if (!$(e.relatedTarget).length) return
    resetModal()
})

$(document).ready(function () {
    loadTemplates()
})

function escapeHtml(s) {
    return String(s)
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
}
