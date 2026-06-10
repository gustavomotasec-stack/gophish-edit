var phoneLists = []
var deleteId = 0

function loadLists() {
    query("/phone_lists/", "GET", {}, true)
    .done(function (data) {
        phoneLists = data || []
        renderTable()
    })
}

function renderTable() {
    var tbody = $("#listsTable tbody").empty()
    if (!phoneLists.length) {
        $("#emptyMessage").show()
        $("#listsTable").hide()
        return
    }
    $("#emptyMessage").hide()
    $("#listsTable").show()
    $.each(phoneLists, function (i, p) {
        var date = moment(p.modified_date).format("DD/MM/YYYY HH:mm")
        var count = p.numbers ? p.numbers.length : 0
        var row = '<tr>' +
            '<td>' + escapeHtml(p.name) + '</td>' +
            '<td><span class="badge">' + count + '</span></td>' +
            '<td>' + date + '</td>' +
            '<td>' +
            '<button class="btn btn-xs btn-default" onclick="editList(' + p.id + ')"><i class="fa fa-pencil"></i> Editar</button> ' +
            '<button class="btn btn-xs btn-danger" onclick="confirmDelete(' + p.id + ')"><i class="fa fa-trash"></i> Excluir</button>' +
            '</td></tr>'
        tbody.append(row)
    })
}

function parseNumbers(text) {
    return text.split('\n')
        .map(function (n) { return n.replace(/\D/g, '') })
        .filter(function (n) { return n.length >= 10 })
        .map(function (n) { return { number: n } })
}

function updateCount() {
    var nums = parseNumbers($("#numbers").val())
    $("#count-info").text(nums.length + " número(s) válido(s)")
}

function resetModal() {
    $("#list-id").val("0")
    $("#modal-title").text("Nova Phone List")
    $("#name").val("")
    $("#numbers").val("")
    $("#count-info").text("")
    $("#modal-alert").empty()
}

function editList(id) {
    var p = phoneLists.find(function (x) { return x.id === id })
    if (!p) return
    $("#list-id").val(p.id)
    $("#modal-title").text("Editar Phone List")
    $("#name").val(p.name)
    var nums = (p.numbers || []).map(function (n) { return n.number }).join('\n')
    $("#numbers").val(nums)
    updateCount()
    $("#modal").modal("show")
}

function saveList() {
    var id = parseInt($("#list-id").val())
    var numbers = parseNumbers($("#numbers").val())
    if (!numbers.length) {
        $("#modal-alert").html('<div class="alert alert-danger">Adicione ao menos um número válido.</div>')
        return
    }
    var data = { name: $("#name").val(), numbers: numbers }
    var method = id === 0 ? "POST" : "PUT"
    var endpoint = id === 0 ? "/phone_lists/" : "/phone_lists/" + id

    query(endpoint, method, data, true)
    .done(function (p) {
        if (id === 0) {
            phoneLists.push(p)
        } else {
            phoneLists = phoneLists.map(function (x) { return x.id === id ? p : x })
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
    query("/phone_lists/" + deleteId, "DELETE", {}, true)
    .done(function () {
        phoneLists = phoneLists.filter(function (x) { return x.id !== deleteId })
        renderTable()
        $("#deleteModal").modal("hide")
    })
})

$("#numbers").on("input", updateCount)

$("#modal").on("show.bs.modal", function (e) {
    if (!$(e.relatedTarget).length) return
    resetModal()
})

$(document).ready(function () {
    loadLists()
})

function escapeHtml(s) {
    return String(s)
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
}
