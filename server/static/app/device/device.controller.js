'use strict';

angular.module('aptWebApp')
    .controller('DeviceDetailCtrl', function($scope, $http, $stateParams, $window, $interval) {
        $scope.deviceID = $stateParams.id;
        $scope.nodeIP = '';
        $scope.device = {};
        $scope.logs = [];
        var logs = [];

        //the status of device
        $scope.getState = function(state) {
            if (state == 0) {
                return 'busy';
            } else {
                return 'free';
            }
        };

        //send device button events to websocket
        $scope.deviceButton = function(kind) {
            if ($scope.screenWS) {
                $scope.screenWS.send(kind)
            }
        }

        $http.get('device/ip/' + $scope.deviceID).then(function(response) {
            if (response) {
                $scope.nodeIP = response.data;

                //create websocket
                var wsUrl = 'ws://' + $scope.nodeIP + ':9002/' + $scope.deviceID;
                var ws = new WebSocket(wsUrl);
                $scope.screenWS = ws;
                //ws.binaryType = 'blob';

                ws.onclose = function() {
                    console.log('onclose');
                };
                //screen display
                ws.onerror = function() {
                    console.log('onerror');
                };
                var BLANK_IMG = 'data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==';
                var element = $scope.drawElement;
                var canvas = $scope.drawElement[0];
                var g = canvas.getContext('2d');
                ws.onmessage = function(message) {
                    if (typeof message.data == 'string') {
                        //log content
                        var segments = message.data.split(' ');
                        var log = {};
                        if (segments.length < 4) {
                            log.kind = '';
                            log.date = '';
                            log.content = message.data;
                            logs.push(log);
                        } else {
                            log.kind = segments[2][0];
                            log.date = segments[0] + ' ' + segments[1];
                            log.content = '';
                            for (var i = 2; i < segments.length; i++) {
                                log.content = log.content + ' ' + segments[i]
                            }
                            logs.push(log);
                        }
                    } else {
                        //image binary
                        var blob = new Blob([message.data], {
                            type: 'image/jpeg'
                        });
                        var URL = window.URL || window.webkitURL
                        var img = new Image();
                        img.onload = function() {
                            //console.log(img.width, img.height)
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
                    }
                };

                ws.onopen = function() {
                    console.log('onopen', arguments);
                    ws.send('1920x1080/0');
                };

                //UI events
                var drawing = false;
                var lastX;
                var lastY;
                element.bind('mousedown', function(event) {
                    if (event.offsetX !== undefined) {
                        lastX = event.offsetX;
                        lastY = event.offsetY;
                    } else {
                        lastX = event.layerX - event.currentTarget.offsetLeft;
                        lastY = event.layerY - event.currentTarget.offsetTop;
                    }
                    drawing = true;
                });
                // element.bind('mousemove', function(event) {
                //     console.log(event);
                // });
                var sendme = function sendEvent(x1, y1, x2, y2) {
                    var dist = (x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)
                    if (dist < 4) {
                        var evt = x1 + ',' + y1;
                        ws.send(evt)
                    } else {
                        var evt = x1 + ',' + y1 + ',' + x2 + ',' + y2;
                        ws.send(evt);
                    }
                };

                element.bind('mouseup', function(event) {
                    if (drawing) {
                        var currentX;
                        var currentY;
                        // get current mouse position
                        if (event.offsetX !== undefined) {
                            currentX = event.offsetX;
                            currentY = event.offsetY;
                        } else {
                            currentX = event.layerX - event.currentTarget.offsetLeft;
                            currentY = event.layerY - event.currentTarget.offsetTop;
                        }
                        sendme(lastX, lastY, currentX, currentY);
                        drawing = false;
                    }
                });

            }
        }, function(response) {

        });

        $scope.refresh = function() {
            $http.get('device/' + $scope.deviceID).then(function(response) {
                if (response) {
                    $scope.device = response.data;
                }
            }, function(response) {
                //console.log(response)
            });
        }

        $scope.refresh();

        var interval = $interval(function() {
            $scope.logs = logs;
        }, 1000);

        $scope.$on('$destroy', function() {
            if ($scope.screenWS) {
                $scope.screenWS.close();
            }
            $interval.cancel(interval);
        });
    });
