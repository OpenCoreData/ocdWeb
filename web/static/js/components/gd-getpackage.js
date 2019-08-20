import {
	html,
	render
} from './lit-html.js';

(function () {
    class SimpleGet extends HTMLElement {
        constructor() {
            super();

            const resID = this.getAttribute('res-id');

            // CAUTION DEV / DEMO HACK..  comment out in production!!!!!!
            var newstr = resID.replace(/opencoredata.org/i, '192.168.2.89:9900');
            // var newstr = "http://192.168.2.89:9900/id/do/blavhlau6s74ma5rn4ng"
            // console.log(newstr); 

            // GET test
            function tj_providers(id) {
                return fetch(id, {
                    headers: { 'Accept': 'application/ld+json', },
                    })
                    .then(function (response) {
                        return response.json();
                    })
                    .then(function (myJson) {
                        //console.log(myJson);
                        // console.log(JSON.stringify(myJson));
                        // return JSON.stringify(myJson);
                        return myJson;
                    });
            }

            // GET test call...
            tj_providers(newstr).then((feature) => {
                this.attachShadow({ mode: 'open' });


                this.shadowRoot.innerHTML = `
                <div style="overflow-wrap: break-word;width=100%">
                    Description: ${feature["description"]} <br>
                    Distribution: <a target="_blank" href="${feature["distribution"].contentUrl}">
                    Download Package Description</a> <br>
                </div>
                  `;


                // var count = Object.keys(providers).length;
                // const itemTemplates = [];
                // var i;
                // for (i = 0; i < count; i++) {
                //     // console.log(providers[i].name)
                //     itemTemplates.push(  `${providers[i].name}`);
                //     // console.log(itemTemplates)
                // }

                // var h =  `<div>${itemTemplates}</div>`;
                // this.shadowRoot.innerHTML = `${h}` ;
                // this.shadowRoot.appendChild(this.cloneNode(h));
            });
        }
    }
    window.customElements.define('geodex-getpackage', SimpleGet);
})();


