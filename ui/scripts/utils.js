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
// Event object
//-----------------------------------------------------------

var Event = {
		_listeners: {},
		_eventNumber: 0,

		on: function (events, callback) {
				this._listeners[events + --this._eventNumber] = callback;
		},

		off: function (events) {
			delete this._listeners[events];
		},

		notify: function (events, data) {
			for ( var topic in this._listeners) {
				if (this._listeners.hasOwnProperty(topic)) {
					if (topic.split("-")[0] == events) {
						this._listeners[topic](data) !== false || delete this._listeners[topic];
					}
				}
			}

		} 
};

//
// Model Object
//-----------------------------------------------------------
var Model = function (attributes) {
	this.id = _.uniqueId('model');
	this.attributes = attributes || {};
};

Model.prototype.get = function(attr) {
	return this.attributes[attr];
};

Model.prototype.set = function(attrs){
	if (_.isObject(attrs)) {

		_.extend(this.attributes, attrs);
		this.change(attrs);
	}
	return this;
};

Model.prototype.toJSON = function(options) {
	return _.clone(this.attributes);
};

Model.prototype.change = function(attrs){
	this.notify(this.id + 'update', attrs);
}; 

_.extend(Model.prototype, Event);


//
// View Object
//-----------------------------------------------------------
var View = function (options) {
	_.extend(this, options); 
	this.id = _.uniqueId('view');
};

_.extend(View.prototype, Event);

//
// Controller Object
//-----------------------------------------------------------
var Controller = function(options) {
	_.extend(this, options);
	this.id = _.uniqueId('controller');		
};