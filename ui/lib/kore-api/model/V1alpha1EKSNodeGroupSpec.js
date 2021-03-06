/**
 * Kore API
 * Kore API provides the frontend API (kore.appvia.io)
 *
 * The version of the OpenAPI document: 0.0.1
 * Contact: info@appvia.io
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';
import V1Ownership from './V1Ownership';

/**
 * The V1alpha1EKSNodeGroupSpec model module.
 * @module model/V1alpha1EKSNodeGroupSpec
 * @version 0.0.1
 */
class V1alpha1EKSNodeGroupSpec {
    /**
     * Constructs a new <code>V1alpha1EKSNodeGroupSpec</code>.
     * @alias module:model/V1alpha1EKSNodeGroupSpec
     * @param amiType {String} 
     * @param credentials {module:model/V1Ownership} 
     * @param desiredSize {Number} 
     * @param diskSize {Number} 
     * @param eC2SSHKey {String} 
     * @param enableAutoscaler {Boolean} 
     * @param maxSize {Number} 
     * @param minSize {Number} 
     * @param region {String} 
     * @param subnets {Array.<String>} 
     */
    constructor(amiType, credentials, desiredSize, diskSize, eC2SSHKey, enableAutoscaler, maxSize, minSize, region, subnets) { 
        
        V1alpha1EKSNodeGroupSpec.initialize(this, amiType, credentials, desiredSize, diskSize, eC2SSHKey, enableAutoscaler, maxSize, minSize, region, subnets);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, amiType, credentials, desiredSize, diskSize, eC2SSHKey, enableAutoscaler, maxSize, minSize, region, subnets) { 
        obj['amiType'] = amiType;
        obj['credentials'] = credentials;
        obj['desiredSize'] = desiredSize;
        obj['diskSize'] = diskSize;
        obj['eC2SSHKey'] = eC2SSHKey;
        obj['enableAutoscaler'] = enableAutoscaler;
        obj['maxSize'] = maxSize;
        obj['minSize'] = minSize;
        obj['region'] = region;
        obj['subnets'] = subnets;
    }

    /**
     * Constructs a <code>V1alpha1EKSNodeGroupSpec</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1alpha1EKSNodeGroupSpec} obj Optional instance to populate.
     * @return {module:model/V1alpha1EKSNodeGroupSpec} The populated <code>V1alpha1EKSNodeGroupSpec</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1alpha1EKSNodeGroupSpec();

            if (data.hasOwnProperty('amiType')) {
                obj['amiType'] = ApiClient.convertToType(data['amiType'], 'String');
            }
            if (data.hasOwnProperty('cluster')) {
                obj['cluster'] = V1Ownership.constructFromObject(data['cluster']);
            }
            if (data.hasOwnProperty('credentials')) {
                obj['credentials'] = V1Ownership.constructFromObject(data['credentials']);
            }
            if (data.hasOwnProperty('desiredSize')) {
                obj['desiredSize'] = ApiClient.convertToType(data['desiredSize'], 'Number');
            }
            if (data.hasOwnProperty('diskSize')) {
                obj['diskSize'] = ApiClient.convertToType(data['diskSize'], 'Number');
            }
            if (data.hasOwnProperty('eC2SSHKey')) {
                obj['eC2SSHKey'] = ApiClient.convertToType(data['eC2SSHKey'], 'String');
            }
            if (data.hasOwnProperty('enableAutoscaler')) {
                obj['enableAutoscaler'] = ApiClient.convertToType(data['enableAutoscaler'], 'Boolean');
            }
            if (data.hasOwnProperty('instanceType')) {
                obj['instanceType'] = ApiClient.convertToType(data['instanceType'], 'String');
            }
            if (data.hasOwnProperty('labels')) {
                obj['labels'] = ApiClient.convertToType(data['labels'], {'String': 'String'});
            }
            if (data.hasOwnProperty('maxSize')) {
                obj['maxSize'] = ApiClient.convertToType(data['maxSize'], 'Number');
            }
            if (data.hasOwnProperty('minSize')) {
                obj['minSize'] = ApiClient.convertToType(data['minSize'], 'Number');
            }
            if (data.hasOwnProperty('region')) {
                obj['region'] = ApiClient.convertToType(data['region'], 'String');
            }
            if (data.hasOwnProperty('releaseVersion')) {
                obj['releaseVersion'] = ApiClient.convertToType(data['releaseVersion'], 'String');
            }
            if (data.hasOwnProperty('sshSourceSecurityGroups')) {
                obj['sshSourceSecurityGroups'] = ApiClient.convertToType(data['sshSourceSecurityGroups'], ['String']);
            }
            if (data.hasOwnProperty('subnets')) {
                obj['subnets'] = ApiClient.convertToType(data['subnets'], ['String']);
            }
            if (data.hasOwnProperty('tags')) {
                obj['tags'] = ApiClient.convertToType(data['tags'], {'String': 'String'});
            }
            if (data.hasOwnProperty('version')) {
                obj['version'] = ApiClient.convertToType(data['version'], 'String');
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getAmiType() {
        return this.amiType;
    }

    /**
     * @param {String} amiType
     */
    setAmiType(amiType) {
        this['amiType'] = amiType;
    }
/**
     * @return {module:model/V1Ownership}
     */
    getCluster() {
        return this.cluster;
    }

    /**
     * @param {module:model/V1Ownership} cluster
     */
    setCluster(cluster) {
        this['cluster'] = cluster;
    }
/**
     * @return {module:model/V1Ownership}
     */
    getCredentials() {
        return this.credentials;
    }

    /**
     * @param {module:model/V1Ownership} credentials
     */
    setCredentials(credentials) {
        this['credentials'] = credentials;
    }
/**
     * @return {Number}
     */
    getDesiredSize() {
        return this.desiredSize;
    }

    /**
     * @param {Number} desiredSize
     */
    setDesiredSize(desiredSize) {
        this['desiredSize'] = desiredSize;
    }
/**
     * @return {Number}
     */
    getDiskSize() {
        return this.diskSize;
    }

    /**
     * @param {Number} diskSize
     */
    setDiskSize(diskSize) {
        this['diskSize'] = diskSize;
    }
/**
     * @return {String}
     */
    getEC2SSHKey() {
        return this.eC2SSHKey;
    }

    /**
     * @param {String} eC2SSHKey
     */
    setEC2SSHKey(eC2SSHKey) {
        this['eC2SSHKey'] = eC2SSHKey;
    }
/**
     * @return {Boolean}
     */
    getEnableAutoscaler() {
        return this.enableAutoscaler;
    }

    /**
     * @param {Boolean} enableAutoscaler
     */
    setEnableAutoscaler(enableAutoscaler) {
        this['enableAutoscaler'] = enableAutoscaler;
    }
/**
     * @return {String}
     */
    getInstanceType() {
        return this.instanceType;
    }

    /**
     * @param {String} instanceType
     */
    setInstanceType(instanceType) {
        this['instanceType'] = instanceType;
    }
/**
     * @return {Object.<String, String>}
     */
    getLabels() {
        return this.labels;
    }

    /**
     * @param {Object.<String, String>} labels
     */
    setLabels(labels) {
        this['labels'] = labels;
    }
/**
     * @return {Number}
     */
    getMaxSize() {
        return this.maxSize;
    }

    /**
     * @param {Number} maxSize
     */
    setMaxSize(maxSize) {
        this['maxSize'] = maxSize;
    }
/**
     * @return {Number}
     */
    getMinSize() {
        return this.minSize;
    }

    /**
     * @param {Number} minSize
     */
    setMinSize(minSize) {
        this['minSize'] = minSize;
    }
/**
     * @return {String}
     */
    getRegion() {
        return this.region;
    }

    /**
     * @param {String} region
     */
    setRegion(region) {
        this['region'] = region;
    }
/**
     * @return {String}
     */
    getReleaseVersion() {
        return this.releaseVersion;
    }

    /**
     * @param {String} releaseVersion
     */
    setReleaseVersion(releaseVersion) {
        this['releaseVersion'] = releaseVersion;
    }
/**
     * @return {Array.<String>}
     */
    getSshSourceSecurityGroups() {
        return this.sshSourceSecurityGroups;
    }

    /**
     * @param {Array.<String>} sshSourceSecurityGroups
     */
    setSshSourceSecurityGroups(sshSourceSecurityGroups) {
        this['sshSourceSecurityGroups'] = sshSourceSecurityGroups;
    }
/**
     * @return {Array.<String>}
     */
    getSubnets() {
        return this.subnets;
    }

    /**
     * @param {Array.<String>} subnets
     */
    setSubnets(subnets) {
        this['subnets'] = subnets;
    }
/**
     * @return {Object.<String, String>}
     */
    getTags() {
        return this.tags;
    }

    /**
     * @param {Object.<String, String>} tags
     */
    setTags(tags) {
        this['tags'] = tags;
    }
/**
     * @return {String}
     */
    getVersion() {
        return this.version;
    }

    /**
     * @param {String} version
     */
    setVersion(version) {
        this['version'] = version;
    }

}

/**
 * @member {String} amiType
 */
V1alpha1EKSNodeGroupSpec.prototype['amiType'] = undefined;

/**
 * @member {module:model/V1Ownership} cluster
 */
V1alpha1EKSNodeGroupSpec.prototype['cluster'] = undefined;

/**
 * @member {module:model/V1Ownership} credentials
 */
V1alpha1EKSNodeGroupSpec.prototype['credentials'] = undefined;

/**
 * @member {Number} desiredSize
 */
V1alpha1EKSNodeGroupSpec.prototype['desiredSize'] = undefined;

/**
 * @member {Number} diskSize
 */
V1alpha1EKSNodeGroupSpec.prototype['diskSize'] = undefined;

/**
 * @member {String} eC2SSHKey
 */
V1alpha1EKSNodeGroupSpec.prototype['eC2SSHKey'] = undefined;

/**
 * @member {Boolean} enableAutoscaler
 */
V1alpha1EKSNodeGroupSpec.prototype['enableAutoscaler'] = undefined;

/**
 * @member {String} instanceType
 */
V1alpha1EKSNodeGroupSpec.prototype['instanceType'] = undefined;

/**
 * @member {Object.<String, String>} labels
 */
V1alpha1EKSNodeGroupSpec.prototype['labels'] = undefined;

/**
 * @member {Number} maxSize
 */
V1alpha1EKSNodeGroupSpec.prototype['maxSize'] = undefined;

/**
 * @member {Number} minSize
 */
V1alpha1EKSNodeGroupSpec.prototype['minSize'] = undefined;

/**
 * @member {String} region
 */
V1alpha1EKSNodeGroupSpec.prototype['region'] = undefined;

/**
 * @member {String} releaseVersion
 */
V1alpha1EKSNodeGroupSpec.prototype['releaseVersion'] = undefined;

/**
 * @member {Array.<String>} sshSourceSecurityGroups
 */
V1alpha1EKSNodeGroupSpec.prototype['sshSourceSecurityGroups'] = undefined;

/**
 * @member {Array.<String>} subnets
 */
V1alpha1EKSNodeGroupSpec.prototype['subnets'] = undefined;

/**
 * @member {Object.<String, String>} tags
 */
V1alpha1EKSNodeGroupSpec.prototype['tags'] = undefined;

/**
 * @member {String} version
 */
V1alpha1EKSNodeGroupSpec.prototype['version'] = undefined;






export default V1alpha1EKSNodeGroupSpec;

