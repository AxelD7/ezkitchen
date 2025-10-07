function setupAccordion(className, displayType) {
    let items = document.getElementsByClassName(className);
    for (let i = 0; i < items.length; i++) {
        items[i].addEventListener("click", function () {
            this.classList.toggle("active");
            if (this.nextElementSibling.style.display == "none") {
                this.nextElementSibling.style.display = displayType;
            } else {
                this.nextElementSibling.style.display = "none"
            }
        });
    }
}

setupAccordion("accordion", "block");
setupAccordion("subcategory-accordion", "flex");


let addProductBtn = document.getElementsByClassName("add-product-btn")
for (let i = 0; i < addProductBtn.length; i++) {

    addProductBtn[i].addEventListener("click", async function () {
        let category = this.id
        let subcat = this.closest(".subcategory-panel").id
        let color = ""
        let url = `http://localhost:4000/product/get/?category=${category}&subcategory=${subcat}&color=${color}`
        console.log(addProductBtn[i].id)
        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`Response Status: ${response.status}`)
            }
            const modalHTML = await response.text()
            console.log(category)
            document.querySelector(".product-modal").innerHTML = modalHTML
            document.querySelector(".product-modal").showModal()

            let addItemBtn = document.getElementsByClassName('add-item-btn')
            for (let i = 0; i < addItemBtn.length; i++) {
                addItemBtn[i].addEventListener("click", function () {
                    console.log(`I've been pressed for the item! with id`)
                    console.log(addItemBtn[i].id)
                })
            }

        } catch (error) {
            console.error(error.message)
        }

    })
}

