import {
	html,
	render
} from './lit-html.js';



const activesearch = () => {
	console.log("There is an active search underway..  waiting for response")

	return html`
	<div class="loader">Loading...</div>
		`;
}


// lit-html constant
const nusearch = (barval, q) => {
	console.log("nusearchtemplate-----------------------------------------------")
	var count = Object.keys(barval.Results.Bindings).length;

	// At this point we need to see if count is 0 and do something about it.
	if (count < 1) {
		return html`<div style="text-align:center;margin-top:50px;position:relative"> 
			<img style="margin-left:40px;height:40px" src="./images/empty.svg">
			<br>  Sorry, in the scope of items indexed there are no results.</div>`;
	}

	const itemTemplates = [];

	var i;
	for (i = 0; i < count; i++) {
		var s = `${barval.Results.Bindings[i].s.Value}`;
		var g = `${barval.Results.Bindings[i].g.Value}`;
		var type = `${barval.Results.Bindings[i].type.Value}`;
		var lit = `${barval.Results.Bindings[i].literal.Value}`;
		var score = `${barval.Results.Bindings[i].score.Value}`;


		// var locationname = `${barval.Results.Bindings[i].locationname.Value}`;
		// var proj = `${barval.Results.Bindings[i].proj.Value}`;
		// var pi = `${barval.Results.Bindings[i].pi.Value}`;
		// var country = `${barval.Results.Bindings[i].country.Value}`;
		// var state_province = `${barval.Results.Bindings[i].state_province.Value}`;
		// var lat = `${barval.Results.Bindings[i].lat.Value}`;
		// var long = `${barval.Results.Bindings[i].long.Value}`;
		// var score = `${barval.Results.Bindings[i].score.Value}`;

		// KEEP..  orignalMain Item div template
		// 	itemTemplates.push(html `<div class="resultitem" style="margin-top:15px">
		//     Project <a target="_blank" href="/collections/csdco/project/${proj}">${proj}</a> (<a href="./csdco.html?q=${proj}"><img style="height:15px" src="/images/reflect.png"></a>)
		// 	<br>
		// 	A project at ${locationname} 
		// 	(<a href="./csdco.html?q=${locationname}"><img style="height:15px" src="/images/reflect.png"></a>)
		// 	by: ${pi} location: ${state_province}, ${country}
		// 	<br/> 
		// 	Associated hole ID:<br/> <a target="_blank" href="${rurl}">${rurl}</a> 
		// 		<br/>
		//      <span>Spatial coodinates: ${lat}  ${long}... </span>
		// 	 <br/>

		// 	<span style="font-size: smaller;" >(${score}  ) <span> </div>`);
		// }

		itemTemplates.push(html`<div class="resultitem" style="margin-top:15px">

		<br/> <a target="_blank" href="${s}">${s}</a> a ${type}
			<br>
			<details>
   			 	<summary>${lit.substring(0,200)}...</summary>
    			${lit}
			</details>
			
			${score} in ${g}
	      </div>`);
	}

	return html`
	  <div>
		   ${itemTemplates}
      </div>
		`;
};

const query1 = (q, n, s) => {
	return `
	{
		"search_request": {
		  "query": {
			"query": "${q}"
		  },
		  "size": ${n},
		  "from": ${s},
		  "fields": [
			"*"
		  ],
		  "sort": [
			"-_score"
		  ],
		  "highlight": {
			"style": "html",
			"fields": [
			  "name",
			  "description"
			]
		  }
		}
	  }
	  `
}


// popstate for history button
window.onpopstate = event => {
	console.log("opnpopstate seen")
	console.log(event.state)
	//window.location.reload()
}


// core init code
let params = (new URL(location)).searchParams;
let q = params.get('q');
let n = params.get('n');
let s = params.get('s');
let i = params.get('i');

// trap n = null to prime the number return do
if (n == null) {
	n = 20
}

// trap s = nul and prime to 0
if (s == null) {
	s = 0
}

// Set the values of the query boxes based on URL parameters
let qdo = document.querySelector('#q');
let ndo = document.querySelector('#nn');
let sdo = document.querySelector('#s');
let ido = document.querySelector('#i');
qdo.value = q;
ndo.value = n;
sdo.value = s;
ido.value = i;

// if q is not null..   fire off a search, 
if (q != null) {
	searchActions();
}


// event listeners
document.querySelector('#q').addEventListener('keyup', function (e) {
	if (e.keyCode === 13) {
		searchActions();
	}
});

document.querySelector('#update').addEventListener('click', searchActions);
//document.querySelector('#providers').addEventListener('click', providerList);


// --------  funcs and constants below here   ---------------------
function searchActions() {
	// let params = (new URL(location)).searchParams;
	let q = document.querySelector('#q').value
	let s = document.querySelector('#s').value
	let n = document.querySelector('#nn').value
	// let s = params.get('s');
	// let i = params.get('i');

	updateURL();

	// Different search options
	blastsearchsimple(q, n, s);
	// threadSearch(q, n, s, i); 
	// simpleSearch();

	// updateNav();   // write to content div 1
}

// RENAME..   now a REST call to a proxy to SPARQL
function blastsearchsimple(q, n, s) {
	// var formData = new FormData();
	//var data = query1(q, n, s);
	//console.log(data)

	// put up a search being done notification
	const el = document.querySelector('#container2');
	render(activesearch(), el)

	//fetch(`http://geodex.org/api/v1/textindex/getnusearch?q=${data}`)
	// fetch(`https://opencoredata.org/api/beta/graph/csdco/search?q=${q}`)
	fetch(`http://localhost:6789/api/beta/graph/csdco/search?q=${q}`)
		.then(function (response) {
			return response.json();
		})
		.then(function (myJson) {
			console.log(myJson);
			const el = document.querySelector('#container2');
			const navel = document.querySelector('#container1');
			render(nusearch(myJson, q), el);
			// render(navui(myJson.search_result.total_hits), navel);
		});
}

function updateURL() {
	let qdo = document.querySelector('#q');
	let ndo = document.querySelector('#nn');
	let sdo = document.querySelector('#s');
	let ido = document.querySelector('#i');

	let params = new URLSearchParams(location.search.slice(1));
	params.set('q', qdo.value);
	params.set('n', ndo.value);
	params.set('s', sdo.value);
	params.set('i', ido.value);

	//window.history.replaceState({}, '', location.pathname + '?' + params);
	const state = {
		geodexsearch: q
	}
	window.history.pushState({}, '', location.pathname + '?' + params);
}

