clear all; close all; clc;
% addpath (genpath('mlib')) 
% Read and parse db file 
dbFile = './output.json';

db = jsondecode(fileread(dbFile));
fn = fieldnames(db);
subKeys = string(fn);
measurements(length(fn)) = struct();
for i=1:numel(fn)
    measurements(i).SubjID = db.(fn{i}).subj_id;
    
    keys = fieldnames(db.(fn{i}).data_num);
    if isempty(keys)
        measurements(i).data = containers.Map('KeyType','char','ValueType','double');
        continue
    end
    values = zeros(length(keys),1);
    for j = 1:numel(keys)
        values(j)= db.(fn{i}).data_num.(keys{j});
    end
    measurements(i).data = containers.Map(keys, values);
end

keys = ["SHOULDER_WIDTH_X12","SHOULDER_CIRCUMFERENCE_X17","HEAD_WIDTH_X1","HEAD_DEPTH_X3","HEAD_HEIGHT_X2"];
values = [460,1070,190,220,210];
r = MinDist(measurements, keys, values);
r


function  r = MinDist(measurements, keyL, values)
    keyMax = containers.Map('KeyType','char','ValueType','double');
    keyMin = containers.Map('KeyType','char','ValueType','double');

    for k = 1:length(keyL)
        key = convertStringsToChars(keyL(k));
        keyFound = false;
        keyMin(key) = realmin('double');
        keyMax(key) = realmax('double');
        
        for i = 1:length(measurements)
            for tKey = keys(measurements(i).data)
                if strcmp (key,tKey{1})
                    keyFound = true;
                    tVal = measurements(i).data(tKey{1});
                    if tVal > keyMax(key)
                        keyMax(key) = tVal;
                        continue
                    end
                    if tVal> keyMin(key)
                        keyMin(key) = tVal;
                    end
                end
            end
        end
        if keyFound == false
            disp(['Key', key,' not present in any participant'])
        end
    end

    subDist = containers.Map('KeyType','char','ValueType','double');
    for i = 1:length(keyL)
        key = convertStringsToChars(keyL(i));
        weight = 1/(keyMax(key)-keyMin(key));
        for j = 1:length(measurements)
            if isKey(subDist, measurements(j).SubjID) ~= 1
                subDist(measurements(j).SubjID) = 0;
            end
            if isKey (measurements(j).data, key) ~= 1
                subDist(measurements(j).SubjID) = 100000000;
                continue
            end
            subDist(measurements(j).SubjID) = subDist(measurements(j).SubjID) + abs(measurements(j).data(key) - values(i))*weight;
        end
    end
    
    r(length(subDist)) = struct();
    idx =1;
    for k = keys(subDist) 
        r(idx).subj_id = k{1};
        r(idx).dist = subDist (k{1});
        idx=idx+1;
    end
end