{{ define "scripts" }}

	$(document).ready(function(){

		$(".clusternode").click(function () {
			var checkUrl
			var baseUrl = $(this).attr('url')
			var serviceUrl = $(this).attr('url')+"/api/services"

			$.getJSON(serviceUrl, function(services) {
				var servicesinfo = '<div class="servicesdetail">';
				$.each(services.services, function(key, service) {
					servicesinfo += '<div id="'+service.name+'" class="servicedetail">';
						servicesinfo += '<div class="servicename">';
							servicesinfo += service.name;
						servicesinfo += '</div>';					
						
						servicesinfo += '<div id="'+service.name+'" class="servicechecks">';
							for ( c in service.checks) {
								servicesinfo += '<div class="servicecheck" url="'+baseUrl+'" name="'+service.checks[c]+'" status="">'+service.checks[c]+'</div>'
							}
							servicesinfo += '<div style="clear:both;"></div>';
							servicesinfo += '<div class="servicecheckdetails"></div>';
							servicesinfo += '<div class="servicechecksampledetails"></div>';
							servicesinfo += '<div style="clear:both;"></div>';
						servicesinfo += '</div>';
					servicesinfo += '</div>';			
				});
				servicesinfo += '</div>';

				$(".clusternodedetails").html(servicesinfo);

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
			
			$.getJSON(checkUrl, function(check) {
				var checkinfo = '<div class="checkinfo" id="checkname">'+check.name+'</div>';
				checkinfo += '<div class="checkinfo" id="checkdescription">'+check.description+'</div>';
				checkinfo += '<div class="checkinfo" id="checkcommand">'+check.command+'</div>';
				$("div[id='"+parent+"'] .servicecheckdetails").html(checkinfo);
			//end getJSON
			});

			$.getJSON(sampleUrl, function(sample) {
				var sampleinfo = '<div class="sampleinfo" id="sampletime">'+sample.sampletime+'</div>';
				sampleinfo += '<div class="sampleinfo" id="elapsedtime">'+sample.elapsedtime+'</div>';
				sampleinfo += '<div class="sampleinfo" id="expirationtime">'+sample.expirationtime+'</div>';
				sampleinfo += '<div class="sampleinfo" id="output">'+sample.output+'</div>';
				sampleinfo += '<div class="sampleinfo" id="timestamp">'+sample.timestamp+'</div>';
				$("div[id='"+parent+"'] .servicechecksampledetails").html(sampleinfo);
			//end getJSON
			});
		//end click
		});

	}

{{end}}