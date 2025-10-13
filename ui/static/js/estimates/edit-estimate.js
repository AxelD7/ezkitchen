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
            const response = await fetch(`/estimate/${estimateID}/items/`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    product_id: parseInt(this.id),
                    quantity: 1,
                }),
            })
            if (!response.ok) {
                throw new Error(`Response Status: ${response.status}`)
            }
            console.log(addItemBtn[i].id)
            location.reload()
        })
    }
}

function setupDeleteBtn() {
    let delProdBtn = document.getElementsByClassName("delete-prod-btn")
    for (let i = 0; i < delProdBtn.length; i++) {
        delProdBtn[i].addEventListener("click", async function () {
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
setupDeleteBtn()
