'use strict';

angular.module('aptWebApp')
    .directive("drawScreen", function($window, $stateParams) {
        return {
            restrict: "A",
            link: function(scope, element) {
                var BLANK_IMG = 'data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==';
                var canvas = element[0];
                var g = canvas.getContext('2d');
                var deviceID = $stateParams.id;
                var nodeIP = $stateParams.ip;
                var wsUrl = 'ws://' + nodeIP + ':9002/' + deviceID;
                var ws = new WebSocket(wsUrl);
                ws.binaryType = 'blob';

                // ws.onclose = function() {
                //     console.log('onclose');
                // };

                ws.onerror = function() {
                    console.log('onerror');
                };

                ws.onmessage = function(message) {
                    var blob = new Blob([message.data], {
                        type: 'image/jpeg'
                    });
                    var URL = window.URL || window.webkitURL
                    var img = new Image();
                    img.onload = function() {
                        console.log(img.width, img.height)
                        canvas.width = img.width
                        canvas.height = img.height
                        g.drawImage(img, 0, 0)
                        img.onload = null
                        img.src = BLANK_IMG
                        img = null
                        u = null
                        blob = null
                    };
                    var u = URL.createObjectURL(blob);
                    img.src = u;
                };

                ws.onopen = function() {
                    console.log('onopen', arguments);
                    ws.send('1920x1080/0');
                };

                scope.$on('$destroy', function() {
                    ws.close();
                });
            }
        };
    });
