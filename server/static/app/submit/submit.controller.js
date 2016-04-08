'use strict';

angular.module('aptWebApp')
    .controller('SubmitCtrl', function($scope, $http, Upload, $state) {

        //init
        $scope.framework = {
            name: 'monkey'
        };
        $scope.devices = [];
        $scope.monkey = {};
        $scope.robotium = {};

        $scope.frameErr = false;
        $scope.selectorErr = false;
        var idList = [];

        $scope.select = {
            checkAll: false
        };

        //get devices from server
        $http.get('device').then(response => {
            if (response) {
                $scope.devices = response.data;
                for (var i = 0; i < $scope.devices.length; i++) {
                    $scope.devices[i].check = false;
                }
            }
        });

        //check if this form is not completed
        $scope.checkEmpty = function() {

            //get all selected devices
            var devices = $scope.devices;
            idList = [];
            for (var i = 0; i < devices.length; i++) {
                if (devices[i].check) {
                    idList.push(devices[i].Id);
                }
            }

            $scope.frameErr = false;
            $scope.selectorErr = false;
            //check empty input
            if ($scope.framework.name == 'monkey' && (!$scope.monkey.file || !$scope.monkey.pkg || !$scope.monkey.arg)) {
                $scope.frameErr = true;
            } else if ($scope.framework.name == 'robotium' && (!$scope.robotium.app || !$scope.robotium.test)) {
                $scope.frameErr = true;
            } else if (idList.length <= 0) {
                $scope.selectorErr = true;
            }
        }

        //submit this job
        $scope.submitJob = function() {

            $scope.checkEmpty();
            if ($scope.frameErr || $scope.selectorErr) {
                return;
            }

            //create SubJob struct
            var SubJob = {};
            var Data = {};
            if ($scope.framework.name == 'monkey') {
                SubJob.FrameKind = 'monkey';
                SubJob.Frame = {};
                SubJob.Frame.AppPath = $scope.monkey.file.name;
                SubJob.Frame.PkgName = $scope.monkey.pkg;
                SubJob.Frame.Argu = $scope.monkey.arg;
                Data.file = $scope.monkey.file;
            } else if ($scope.framework.name == 'robotium') {
                SubJob.FrameKind = 'robotium';
                SubJob.Frame = {};
                SubJob.Frame.AppPath = $scope.robotium.app.name;
                SubJob.Frame.TestPath = $scope.robotium.test.name;
                Data.app = $scope.robotium.app;
                Data.test = $scope.robotium.test;
            }
            SubJob.FilterKind = 'specify_devices';
            SubJob.Filter = {};
            SubJob.Filter.IdList = idList;
            var jobjson = JSON.stringify(SubJob);
            Data.job = jobjson;

            //submit requirement and files to server
            Upload.upload({
                url: 'job/',
                method: 'POST',
                data: Data
            }).then(function(resp) {
                $state.go('main.jobs');
            }, function(resp) {
                console.log('Error status: ' + resp.status);
            });
        }

        $scope.clickTopBox = function() {
            angular.forEach($scope.devices, function(device) {
                device.check = $scope.select.checkAll;
            });
        }

    });
