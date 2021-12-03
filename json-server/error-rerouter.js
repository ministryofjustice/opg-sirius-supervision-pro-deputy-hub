const getFailRoute = (req) => {
    return req.headers?.cookie?.match(/fail-route=(?<failRoute>.+);/)?.groups
        .failRoute;
};

module.exports = (req, res, next) => {
    if (["POST", "PUT"].includes(req.method)) {
        const failRoute = getFailRoute(req);

        if (failRoute) {
            req.method = "GET";
            req.url = `/errors/${failRoute}`;
            res.status(400);
        }
    }
    next();
};
