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

            // GET test
            function tj_providers(id) {
                return fetch(id, {
                    headers: { 'Content-Type': 'application/ld+json', },
                })
                    .then(function (response) {
                        return response.json();
                    })
                    .then(function (myJson) {
                        //  console.log(myJson);
                        // console.log(JSON.stringify(myJson));
                        // return JSON.stringify(myJson);
                        return myJson;
                    });
            }



            // GET test call...
            tj_providers(newstr).then((feature) => {
                this.attachShadow({ mode: 'open' });

                // CAUTION DEV / DEMO HACK..  comment out in production!!!!!!
                var newid = feature["@id"].replace(/http:\/\/opencoredata.org\/id\/do/i, '.');

                this.shadowRoot.innerHTML = `
                <div style="overflow-wrap: break-word;width=100%">
                    Feature: <a href="${newid}"> 
                    ${feature["http://opencoredata.org/voc/csdco/v1/hole_ID"]}</a>
                    (IGSN: <a href="http://sesar.org/${feature["http://opencoredata.org/voc/csdco/v1/IGSN"]}">
                    ${feature["http://opencoredata.org/voc/csdco/v1/IGSN"]}<a/> 
                    )
                    
                    <br>
                    PI(s): ${feature["http://opencoredata.org/voc/csdco/v1/pi"]}<br>
                     ${feature["http://opencoredata.org/voc/csdco/v1/country"]} > 
                     ${feature["http://opencoredata.org/voc/csdco/v1/county_Region"]} > 
                     ${feature["http://opencoredata.org/voc/csdco/v1/location"]}
                     <br>
                     <a target="_blank" href="https://www.google.com/maps/search/?api=1&zoom=4&basemap=terrain&query=${feature["http://www.w3.org/2003/01/geo/wgs84_pos#lat"]},${feature["http://www.w3.org/2003/01/geo/wgs84_pos#long"]}">
                     (lat:  ${feature["http://www.w3.org/2003/01/geo/wgs84_pos#lat"]}
                      long:  ${feature["http://www.w3.org/2003/01/geo/wgs84_pos#long"]}
                    )
                    </a>
                </div> `;

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
    window.customElements.define('geodex-getfeature', SimpleGet);
})();

