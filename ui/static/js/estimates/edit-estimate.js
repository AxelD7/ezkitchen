function setupAccordion(className, displayType) {
    let items = document.getElementsByClassName(className)
    for (let i = 0; i < items.length; i++) {
        items[i].addEventListener("click", function () {
            this.classList.toggle("active")
            if (this.nextElementSibling.style.display == "none") {
                this.nextElementSibling.style.display = displayType
            } else {
                this.nextElementSibling.style.display = "none"
            }
        })
    }
}

function setupAddProductBtn() {
    let addProductBtn = document.getElementsByClassName("add-product-btn")
    for (let i = 0; i < addProductBtn.length; i++) {
        addProductBtn[i].addEventListener("click", async function () {
            let category = this.id
            let subcat = this.closest(".subcategory-panel").id
            let color = ""
            if (category == subcat) {
                subcat = ""
            }
            let url = `/product/get/?category=${category}&subcategory=${subcat}&color=${color}`
            console.log(addProductBtn[i].id)
            try {
                const response = await fetch(url)
                if (!response.ok) {
                    throw new Error(`Response Status: ${response.status}`)
                }
                const modalHTML = await response.text()
                console.log(category)
                document.querySelector(".product-modal").innerHTML = modalHTML
                document.querySelector(".product-modal").showModal()

                setupAddProductToEstimateBtn()
            } catch (error) {
                console.error(error.message)
            }
        })
    }
}

function setupAddProductToEstimateBtn() {
    let addItemBtn = document.getElementsByClassName("add-item-btn")
    for (let i = 0; i < addItemBtn.length; i++) {
        addItemBtn[i].addEventListener("click", async function () {
            estimateID =
                document.querySelector(".estimate-ID").dataset.estimateId
            quantity = this.closest(".card").querySelector(
                ".card-product-quantity",
            ).value
            const response = await fetch(`/estimate/${estimateID}/items/`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    product_id: parseInt(this.id),
                    quantity: parseInt(quantity),
                }),
            })
            if (!response.ok) {
                throw new Error(`Response Status: ${response.status}`)
            }
            location.reload()
        })
    }
}

function setupUpdateItemBtn() {
    let saveItemBtn = document.getElementsByClassName("update-item-btn")

    for (let i = 0; i < saveItemBtn.length; i++) {
        saveItemBtn[i].addEventListener("click", async function () {
            currentRow = this.closest("tr")
            lineItemID = currentRow.dataset.lineItemId
            newQuantity = currentRow.querySelector(
                ".estimate-item-quantity",
            ).value

            try {
                const response = await fetch(`/estimate/items/${lineItemID}`, {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        quantity: parseInt(newQuantity),
                    }),
                })
                if (!response.ok) {
                    throw new Error(`Response status: ${response.status}`)
                }
                location.reload()
            } catch {
                console.Error(error.message)
            }
        })
    }
}

function setupDeleteBtn() {
    let delItemBtn = document.getElementsByClassName("delete-item-btn")
    for (let i = 0; i < delItemBtn.length; i++) {
        delItemBtn[i].addEventListener("click", async function () {
            lineItemID = this.closest("tr").dataset.lineItemId
            try {
                const response = await fetch(`/estimate/items/${lineItemID}`, {
                    method: "DELETE",
                })
                if (!response.ok) {
                    throw new Error(`Response status: ${response.status}`)
                }
                location.reload()
            } catch {
                console.Error(error.message)
            }
        })
    }
}

setupAccordion("accordion", "block")
setupAccordion("subcategory-accordion", "flex")
setupAddProductBtn()
setupUpdateItemBtn()
setupDeleteBtn()
