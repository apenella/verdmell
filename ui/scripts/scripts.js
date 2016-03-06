{{ define "scripts" }}

	function loadClusterNodeStatus() {
		$(".clusternode").each(function() {
			baseUrl = $(this).attr("url");
			nodeUrl = baseUrl+"/api/node/status";
			
			$.getJSON(nodeUrl, function(node) {
				console.log($(".clusternode[url='"+baseUrl+"']").attr("url"));
				$(".clusternode[url='"+baseUrl+"']").attr("status",node.status);	
			});
		});
	}

	function loadChecksData(url) {
		this.url = url;
		this.data;

		function request(url) { 
			return $.getJSON(url);
		}

		$.when( request(this.url) ).then( function(data){
			this.data = data.checks;

			for (c in this.data){
				console.log(this.data[c].name);
			}
		});
	};

	function loadSamplesData(url) {
		var samples;

		function request(url) { 
			return $.getJSON(url);
		}

		$.when( request(url) ).then( function(data){
			samples = data.Samples;

			for (s in samples) {
				console.log(s+" "+samples[s].Sample.exit);
				$(".servicecheck[name='"+s+"']").attr("status",samples[s].Sample.exit);			
			}
		});
	};

	$(document).ready(function(){
		loadClusterNodeStatus();

		$(".clusternode").click(function () {
			var checkUrl
			var baseUrl = $(this).attr('url')
			var serviceUrl = $(this).attr('url')+"/api/services"
			var checksUrl = $(this).attr('url')+"/api/checks"
			var samplesUrl = $(this).attr('url')+"/api/samples"

			loadChecksData(checksUrl);

			$.getJSON(serviceUrl, function(services) {
				var servicesinfo = '<div class="servicesdetails">';
				$.each(services.services, function(key, service) {
					servicesinfo += '<div id="'+service.name+'" class="servicedetail" status="'+service.status+'">';
						servicesinfo += '<div class="servicename">';
							servicesinfo += service.name;
						servicesinfo += '</div>';					
						
						servicesinfo += '<div id="'+service.name+'" class="servicechecks">';
							for ( c in service.checks) {
								servicesinfo += '<div class="servicecheck" url="'+baseUrl+'" name="'+service.checks[c]+'" status="">'+service.checks[c]+'</div>'
							}
							servicesinfo += '<div style="clear:both;"></div>';
							
							servicesinfo += '<div id="'+service.name+'" class="servicecheckalldetails" check="">';
								servicesinfo += '<div class="servicecheckdetails"></div>';
								servicesinfo += '<div class="servicechecksampledetails"></div>';
							servicesinfo += '</div>';
							servicesinfo += '<div style="clear:both;"></div>';
						servicesinfo += '</div>';
					servicesinfo += '</div>';			
				//end each
				});
				servicesinfo += '</div>';

				$(".clusternodedetails").html(servicesinfo);

				loadSamplesData(samplesUrl);
				setActionToServicesChecks();

			//end getJSON
			});
		//end click
		});


	});


	function setActionToServicesChecks () {

		$(".servicecheck").click(function(){
			var baseUrl = $(this).attr('url');
			var check = $(this).attr('name');
			var checkUrl = baseUrl+"/api/checks/"+check;
			var sampleUrl = baseUrl+"/api/samples/"+check;

			var parent = $(this).parent().attr("id");


			if ( $(this).attr("name") == $("div[id='"+parent+"'] .servicecheckalldetails").attr("check") ){
				$("div[id='"+parent+"'] .servicecheckalldetails").toggle(200);
			} else {
				$("div[id='"+parent+"'] .servicecheckalldetails").attr("check",check);
				$.getJSON(checkUrl, function(check) {
					$("div[id='"+parent+"'] .servicecheckdetails").attr("check",check.name);
					var checkinfo = '<div class="checkinfotitle" id="checkname">Check</div>';
					checkinfo +=  '<div class="checkinfo" id="checkname">'+check.name+'</div>';
					checkinfo += '<div class="checkinfotitle" id="checkname">Description</div>';
					checkinfo += '<div class="checkinfo" id="checkdescription">'+check.description+'</div>';
					checkinfo += '<div class="checkinfotitle" id="checkname">Command</div>';
					checkinfo += '<div class="checkinfo" id="checkcommand">'+check.command+'</div>';
					$("div[id='"+parent+"'] .servicecheckalldetails .servicecheckdetails").html(checkinfo);
				//end getJSON
				});

				$.getJSON(sampleUrl, function(sample) {
					var sampleinfo = '<div class="sampleinfotitle" id="sampletime">Sample time</div>';
					sampleinfo += '<div class="sampleinfo" id="sampletime">'+sample.sampletime+'</div>';
					sampleinfo += '<div class="sampleinfotitle" id="sampletime">Exit value</div>';
					sampleinfo += '<div class="sampleinfo" id="elapsedtime">'+sample.exit+'</div>';
					sampleinfo += '<div class="sampleinfotitle" id="sampletime">Elapsed time (ns)</div>';
					sampleinfo += '<div class="sampleinfo" id="elapsedtime">'+sample.elapsedtime+'</div>';
					sampleinfo += '<div class="sampleinfotitle" id="sampletime">Sample expiration time (ns)</div>';
					sampleinfo += '<div class="sampleinfo" id="expirationtime">'+sample.expirationtime+'</div>';
					sampleinfo += '<div class="sampleinfotitle" id="sampletime">Output</div>';
					sampleinfo += '<div class="sampleinfo" id="output">'+sample.output+'</div>';
					sampleinfo += '<div class="sampleinfotitle" id="sampletime">Timestamp</div>';
					sampleinfo += '<div class="sampleinfo" id="timestamp">'+sample.timestamp+'</div>';
					$("div[id='"+parent+"'] .servicecheckalldetails .servicechecksampledetails").html(sampleinfo);
				//end getJSON
				});
			//end if click to same check
			}
			
		//end click
		});

	}

{{end}}