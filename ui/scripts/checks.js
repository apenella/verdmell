/*
	Package 'ui' 
	-server
	-handler
	-router
	-routes

	-html/
	-images/
	-pages/
	-scripts/
	-style/

*/

//
// Model Object
//-----------------------------------------------------------
var checksModel = new Model('checksModel',{

	_nodes: [],
	_checks: [],

	_itemTypeClass: {
		"nodes": "item",
		"services": "detailitem"
	},

	// setNodes: function(data) {
	// 	this._nodes = [],

	// 	_.each(data, function(node){
	// 		checksModel.attributes.addNode(node.name, node.URL);
	// 	});
	// },

	// addNode: function(node, url) {
	// 	this._nodes.push(checksModel.attributes.createNode(node, url));
	// },

	// getNode: function(node) {
	// 	return _.where(checksModel.attributes._nodes, {node: node});
	// },

	// createNode: function (node, url) {
	// 	return { node: node, url: url, selected: false };
	// },

	add: function(name, locator, selected, show, checks, samples) {
		this._checks.push(checksModel.attributes.create(name, locator, selected, show, checks, samples));
	},

	create: function(name, locator, selected, show, checks, samples) {
		return { name: name, locator: locator, selected: selected, show: show, checks: checks, samples: samples};
	},

	clear: function() {
		// console.log('checksModel::clear');
		this._checks = [];
	},

	setChecks: function() {
		// console.log('checksModel::setChecks');
		checksModel.set(this._checks);
	},

	//show
	show: function() {

	},

	select: function(locator) {
		console.log('checksModel::select',locator);
		// review each item
		_.each(this._checks, function(item){
			// set selected
			if ( item.locator == locator) {
				item.selected = true;
			} else {
				item.selected = false;
			}
		});
		// set the changes an notify subscribers
		checksModel.set(this._elements);
	},

	// generate locator
	generateLocator: function(base, check) {
		return base + '/' + check;
	},

	// return selected items
	getSelected: function(){
		return _.where(this._checks, {selected: true});
	},

	//
	//
	// subscriptions
	observe: function(model, f) {
		clusterlistModel.on(model.id, this.id, function(){f(model);}.bind(this));
	},
	// observe details
	observeDetails: function(model) {
		if (detailsModel.attributes.getSelected().length) {
			// console.log('checksModel::observe',detailsModel.attributes.getSelected()[0].content.checks);
			// get node data related to detailitem clicked	
			node = clusterlistModel.attributes.getNode(detailsModel.attributes.getSelected()[0][checksModel.attributes._itemTypeClass[detailsModel.attributes.getSelected()[0].type]]);
			//node = checksModel.attributes.getNode(detailsModel.attributes.getSelected()[0][checksModel.attributes._itemTypeClass[detailsModel.attributes.getSelected()[0].type]]);
			// console.log('checksModel::observe',node);
			if ( node.length ) {
				//checksModel.attributes.retrieveChecks(node[0].url);
				$.getJSON(node[0].data.URL+MainController.nodeuri, function(data){
					// clear checks content
					checksModel.attributes.clear();
					// console.log('checksModel::observe',data.checks);
					$.each(detailsModel.attributes.getSelected()[0].content.checks, function(index, check){
						// console.log('checksModel::observe',data.checks.checks.checks[check]);
						// console.log('checksModel::observe',data.samples.Samples[check]);
						// console.log('checksModel::observe', checksModel.attributes.generateLocator(detailsModel.attributes.getSelected()[0].locator,data.checks.checks.checks[check].name));
						checksModel.attributes.add(data.checks.checks.checks[check].name, checksModel.attributes.generateLocator(detailsModel.attributes.getSelected()[0].locator,data.checks.checks.checks[check].name),false, false, data.checks.checks.checks[check],data.samples.Samples[check]);
					});

					checksModel.attributes.setChecks();
				});
			}
		}	
	}

});

//
// View Object
//-----------------------------------------------------------
var checksView = new View({
	// place to append the items
	parent: ".checks",

	initialize: function() {
		$('<div/>', {
			class: 'checkslist'
		}).appendTo(checksView.parent);
	},

	render: function(model) {
		$(checksView.parent).empty();

			$('<div/>', {
				class: 'checkstitle',
				id: 'title',
				text: 'checks'
			}).appendTo(checksView.parent);

		// console.log('checksView::render',model);
		_.each(model, function(item){
			$('<div/>', {
				class: 'check',
				id: item.locator,
				status: item.samples.Sample.exit,
				text: item.name
			}).appendTo(checksView.parent);
		});

		// set events for details
		checksController.setEvents();
	},

	observe: function(model) {
		checksModel.on(model.id, this.id, function(model){
			// console.log('checksView::observe', model);
			checksView.render(model);
		});
	}

});

//
// Controller Object
//-----------------------------------------------------------
var checksController = new Controller({
	// menu data management
	model: checksModel,
	// menu view management
	view: checksView,

	events: {
		".check::click": "select"
	},

	initialize: function(data) {

		//
		// set listeners
		//
		// subscribe to detailsModel
		//this.model.attributes.observe(detailsModel);
		this.model.attributes.observe(detailsModel,this.model.attributes.observeDetails);
		// subscribe for changes to checksModel
		this.view.observe(this.model);

		this.view.initialize();
		// set nodes to model
		// this.model.attributes.setNodes(data['nodes']);

//		console.log('checksController::initialize','data',data['nodes']);
//		console.log('checksController::initialize','clusterlist',clusterlistModel.attributes.getNodes());
	},

	// on update
	update: function(data) {
		// this.model.attributes.setNodes(data['nodes']);
	},

	// set events to items
	setEvents: function() {
		var parts, selector, eventType;
		if(this.events){
			_.each(this.events, function(method, eventName){
				parts = eventName.split('::');
				// get item from events property
				selector = parts[0];
				// get method from events property
				eventType = parts[1];
				// hold the event to item
				$(selector).on(eventType, this[method]);
			}.bind(this));
		}
	},

	select: function() {
		console.log('checksController::select');
		checksController.model.attributes.select(this.getAttribute('id'));
	}

});
