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
// Detail's Model Object
//-----------------------------------------------------------
var detailsModel = new Model('detailsModel',{
	_elements: [],

	// set elements
	set: function(data) {
		this._elements = [];
		this._locator = locatorModel.attributes.get();
		this._selected = false;

		//
		// iterate all the items from data to generate the clusterlist
		// achive types
		_.each(data, function(contentType, type){
			// achieve item
			_.each(contentType, function(contentItem, item){
				// achieve detail
				_.each(contentItem, function(contentDetailType, detailType){

					// only are required objects where are defined all the details
					if (_.isObject(contentDetailType)){
						// achive moreDetail
						_.each(contentDetailType,function(contentDetailItem, detailItem){
							locator = detailsModel.attributes.generateLocator(type,item,detailItem);
								detailsModel.attributes.add(type, item, detailType, detailItem, locator, this._selected, false, contentItem.URL, contentDetailItem);
						});// end achive detail item
					}// end isObject
				});// end achieve detail type
			});// end achieve item
		});// end achieve types
	
		// set the array of elements
		// detailsModel.set(this._elements);

		// if there was any selected clusterlist item, this is setted again
		if ( _.size(this._locator) > 0 && _.findWhere(this._elements, {locator: this._locator}) != undefined ) {
		 	clusterlistModel.attributes.select(this._locator);
		}

		// set the element to show
		if (clusterlistModel.attributes.getSelected().length > 0 ) {
			detailsModel.attributes.show(clusterlistModel.attributes.getSelected()[0].locator);
		}
	},

	// create detail item
	create: function(type, item, detailtype, detailitem, locator, selected, show, urlbase, content) {
		return {type: type, item: item, detailtype: detailtype, detailitem: detailitem, selected: selected, show: show, url: urlbase, locator: locator, content: content};
	},

	// add an item to _elements
	add: function(type, item, detail, moreDetail, locator, selected, show, urlbase, content) {
		this._elements.push(this.create(type, item, detail, moreDetail, locator, selected, show, urlbase, content));
	},

	// mark items to be showed
	show: function(locator) {
		// console.log('detailsModel::show',locator);
		
		// if exists, get selected item
		selected = '';
		if (clusterlistModel.attributes.getSelected().length) {
			selected = detailsModel.attributes.getSelected()[0];
		}

		// for each detailModel's element
		_.each(this._elements, function(item){
			// console.log('detailsModel::show',item);
			
			// show items which locator base is the same as the selected item on cluster list.
			if ( locator == detailsModel.attributes.getBase(item.locator)) {
				// set item to be showed
				item.show = true;
				// select item if it was already selected
				//if (selected.locator == item.locator ) item.selected = true;
			} else {
				// console.log('detailsModel::show','show false',item);
				// set item for not being showed
				item.show = false;
				// unselect if item was selected
				item.selected = false;
			}
		});

		detailsModel.set(this._elements);
	},

	// mark an item as selected
	select: function(locator) {
		// review each item
		_.each(this._elements, function(item){
			// set selected
			if ( item.locator == locator) {
				console.log('detailsModel::select',locator);
				item.selected = true;
			} else {
				item.selected = false;
			}
		});
		// set the changes an notify subscribers
		detailsModel.set(this._elements);
	},

	// return item marked to be showed
	getShowed: function() {
		return _.where(this._elements, {show: true});
	},

	// return to which type of item belongs the detail
	getBase: function(locator) {
		locatorSplitted = locator.split('/');
		if ( locatorSplitted.length > 1 )
			return '/' + locatorSplitted[1] +'/'+ locatorSplitted[2]
		else
			return null;
	},

	// return the detail's id
	getDetail: function(locator) {
		locatorSplitted = locator.split('/');
		if ( locatorSplitted.length > 2 )
			return locatorSplitted[3]
		else
			return null;
	},
	// return selected items
	getSelected: function(){
		return _.where(this._elements, {selected: true});
	},

	// generateLocator
	generateLocator: function(type,item,detail) {
		return '/'+ type + '/' + item + '/' + detail;
	},

	//
	//
	// subscriptions
	observe: function(model, f) {
		clusterlistModel.on(model.id, this.id, function(){f(model);}.bind(this));
	},
	// observe cluster list
	observeClusterlist: function(model) {
		// console.log('detailsModel::observeClusterlist',selected);
		// if exist, get selected item from clusterlist
		selected = '';
		if (model.attributes.getSelected().length > 0) {
			selected = model.attributes.getSelected()[0].locator;
		}
		detailsModel.attributes.show(selected);
	}

});

//
// Detail's View Object
//-----------------------------------------------------------
var detailsView = new View({
	// place to append the items
	parent: ".clusterlistitemdetails",
	// titles
	_title: {
		"nodes": "services",
		"services": "nodes"
	},

	//
	//
	// subscriptions
	observe: function(model, f) {
		detailsModel.on(model.id, this.id, function(){f(model);}.bind(this));
	},

	render: function(model) {
 		// console.log('detailsView::render', model.attributes.getShowed().length);
 		$(detailsView.parent).empty();

		if (model.attributes.getShowed().length) {	
			// console.log('detailsView::_render', model.attributes.getShowed().length);
			$('<div/>', {
					class: 'itemdetailstitle',
					text: detailsView._title[model.attributes.getShowed()[0].type].toUpperCase()
			}).appendTo(detailsView.parent);

			_.each(model.attributes.getShowed(), function(item){
				//console.log('detailsView::_render', item);
				 $('<div/>', {
						class: 'itemdetails',
						id: item.locator,
						status: item.content.status,
						text: item.detailitem
				}).appendTo(detailsView.parent);
			});

			// set events for details
			detailsController.setEvents();
		}
	}

});

//
// Controller Object
//-----------------------------------------------------------
var detailsController = new Controller({
	// menu data management
	model: detailsModel,
	// menu view management
	view: detailsView,

	events: {
		".itemdetails::click": "selected"
	},

	//
	// menu initializes when document is ready
	initialize: function(data) {
		//console.log('detailsController::initialize');

		//
		// set listeners
		//
		// subscribe view to model
		//this.view.observe(this.model);
		this.view.observe(this.model,this.view.render);
		// subscribe for changes to clusterlist
		this.model.attributes.observe(clusterlistModel, this.model.attributes.observeClusterlist);
		// this.view.observeClusterlist(clusterlistModel);
		// subscribe for changes to menu
		// this.view.observeMenu(menuModel);

		// initialitze data for clusterlist
		detailsModel.attributes.set(data);

		// set events for details
		//detailsController._setEvents();
	},

	_initializeWorker: function() {

	},

	update: function(data) {
		detailsModel.attributes.set(data);
	},


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

	selected: function() {
		// console.log('detailsController::selected',this.getAttribute('id'));
		detailsController.model.attributes.select(this.getAttribute('id'));
	}

});